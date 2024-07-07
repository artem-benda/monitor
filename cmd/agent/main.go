package main

import (
	"context"
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
	pb "github.com/artem-benda/monitor/internal/grpc/mon"
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

	retryController := retry.NewRetryController(errors.ErrNetwork{}, errors.ErrServerTemporary{})

	// Слушаем нужные сигналы от ОС
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	stdMetricsCh := genStdMetrics(ctx.Done())
	psUtilsMetricsCh := genPSUtilsMetrics(ctx.Done())

	metricsCh := fanIn(stdMetricsCh, psUtilsMetricsCh)

	var wg sync.WaitGroup

	var sendMetrics func(metrics map[model.MetricKey]model.MetricValue) error

	if config.UseGRPC {
		client, conn := mustCreateGRPCClient()
		defer conn.Close()
		sendMetrics = createGRPCSendMetrics(client, retryController)
	} else {
		client := mustCreateRestyClient()
		sendMetrics = createRESTSendMetrics(client, retryController)
	}

	for i := 0; i < config.MaxParallelWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sendMetricsWorker(id, sendMetrics, metricsCh)
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

func sendMetricsWorker(id int, sendMetrics func(map[model.MetricKey]model.MetricValue) error, in <-chan map[model.MetricKey]model.MetricValue) {
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
		err := sendMetrics(metrics)

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

func createRESTSendMetrics(client *resty.Client, retryController retry.RetryController) func(metrics map[model.MetricKey]model.MetricValue) error {
	rsaPublicKey := mustParsePublicKey(config.RSAPubKeyBase64)
	return func(metrics map[model.MetricKey]model.MetricValue) error {
		return requests.SendAllMetrics(client, retryController, metrics, []byte(config.Key), rsaPublicKey)
	}
}

func createGRPCSendMetrics(client pb.MonitorServiceClient, retryController retry.RetryController) func(metrics map[model.MetricKey]model.MetricValue) error {
	return func(metrics map[model.MetricKey]model.MetricValue) error {
		dtos := make([]*pb.MetricValue, 0)
		for k, v := range metrics {
			var dto *pb.MetricValue
			if k.Kind == model.GaugeKind {
				dto = &pb.MetricValue{MetricId: k.Name, Value: &pb.MetricValue_Gauge{Gauge: v.Gauge}}
			} else if k.Kind == model.CounterKind {
				dto = &pb.MetricValue{MetricId: k.Name, Value: &pb.MetricValue_Counter{Counter: v.Counter}}
			}

			dtos = append(dtos, dto)
		}
		err := retryController.Run(func() (err error) {
			_, updErr := client.UpdateMetricsBatch(context.Background(), &pb.UpdateMetricsBatchRequest{Metrics: dtos})
			return updErr
		})

		return err
	}
}
