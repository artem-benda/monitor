package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/artem-benda/monitor/internal/client/errors"
	"github.com/artem-benda/monitor/internal/client/requests"
	"github.com/artem-benda/monitor/internal/client/service"
	"github.com/artem-benda/monitor/internal/client/storage"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/model"
	"github.com/artem-benda/monitor/internal/retry"
	"github.com/artem-benda/monitor/internal/version"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

const (
	metricsCacheSize = 10
	addr             = ":8089" // адрес сервера pprof
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	version.PrintVersion(buildVersion, buildDate, buildCommit)
	parseFlags()

	if err := logger.Initialize(config.LogLevel); err != nil {
		panic(err)
	} else {
		defer func() {
			err = logger.Log.Sync()
			if err != nil {
				panic(err)
			}
		}()
	}

	client := resty.New()
	serverEndpointURL := fmt.Sprintf("http://%s", config.ServerEndpoint)

	logger.Log.Debug("Starting with base URL", zap.String("baseURL", serverEndpointURL))

	client.SetBaseURL(serverEndpointURL)
	client.SetTimeout(30 * time.Second)
	client.OnAfterResponse(logger.NewRestyResponseLogger())
	client.Header.Add("X-Real-IP", mustGetLocalIPAddr(config.ServerEndpoint).String())

	retryController := retry.NewRetryController(errors.ErrNetwork{}, errors.ErrServerTemporary{})

	// Слушаем нужные сигналы от ОС
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	stdMetricsCh := genStdMetrics(ctx.Done())
	psUtilsMetricsCh := genPSUtilsMetrics(ctx.Done())

	metricsCh := fanIn(stdMetricsCh, psUtilsMetricsCh)

	var wg sync.WaitGroup

	for i := 0; i < config.MaxParallelWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sendMetricsRESTWorker(id, client, retryController, metricsCh)
		}(i)
	}

	go func() {
		log.Fatal(http.ListenAndServe(addr, nil)) // запускаем сервер pprof
	}()

	wg.Wait()
}

func genStdMetrics(doneCh <-chan struct{}) <-chan map[model.MetricKey]model.MetricValue {
	out := make(chan map[model.MetricKey]model.MetricValue)

	go func() {
		defer close(out)
		for {
			select {
			case <-doneCh:
				return
			default:
				out <- service.ReadMetrics(storage.CounterStore)
				logger.Log.Debug("added std metrics")
				time.Sleep(time.Duration(config.PollInterval) * time.Second)
			}
		}
	}()

	return out
}

func genPSUtilsMetrics(doneCh <-chan struct{}) <-chan map[model.MetricKey]model.MetricValue {
	out := make(chan map[model.MetricKey]model.MetricValue)

	go func() {
		defer close(out)
		for {
			select {
			case <-doneCh:
				return
			default:
				out <- service.ReadPSUtilsMetrics()
				logger.Log.Debug("added PSUtils metrics")
				time.Sleep(time.Duration(config.PollInterval) * time.Second)
			}
		}
	}()

	return out
}

func fanIn(chs ...<-chan map[model.MetricKey]model.MetricValue) <-chan map[model.MetricKey]model.MetricValue {
	out := make(chan map[model.MetricKey]model.MetricValue, metricsCacheSize)
	var wg sync.WaitGroup

	for _, ch := range chs {
		chClosure := ch
		wg.Add(1)

		go func() {
			defer wg.Done()
			for i := range chClosure {
				out <- i
			}
		}()
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

func sendMetricsRESTWorker(id int, client *resty.Client, retryController retry.RetryController, in <-chan map[model.MetricKey]model.MetricValue) {
	rsaPublicKey := mustParsePublicKey(config.RSAPubKeyBase64)
	isShutdown := false
	for {
		metrics := make(map[model.MetricKey]model.MetricValue)
	L:
		for {
			select {
			// Выбираем все накопившиеся метрики и мерджим их
			case m, ok := <-in:
				if ok {
					for k, v := range m {
						metrics[k] = v
					}
				} else {
					// Канал закрыт, отправляем что есть
					isShutdown = true
					break L
				}
			// Больше нечего читать, отправляем что выбрали
			default:
				break L
			}
		}
		logger.Log.Debug("sending metrics", zap.Int("workerId", id))
		err := requests.SendAllMetrics(client, retryController, metrics, []byte(config.Key), rsaPublicKey)

		if err != nil {
			logger.Log.Debug("error sending metrics batch", zap.Int("workerId", id), zap.Error(err))
		}

		if isShutdown {
			logger.Log.Debug("gracefully finishing worker", zap.Int("workerId", id))
			return
		}

		storage.CounterStore.Reset()
		time.Sleep(time.Duration(config.ReportInterval) * time.Second)
	}
}
