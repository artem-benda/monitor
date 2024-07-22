package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/artem-benda/monitor/internal/crypt"
	"github.com/artem-benda/monitor/internal/gzipper"
	"github.com/artem-benda/monitor/internal/ipfilter"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/retry"
	"github.com/artem-benda/monitor/internal/server/errors"
	"github.com/artem-benda/monitor/internal/server/handlers"
	"github.com/artem-benda/monitor/internal/server/storage"
	"github.com/artem-benda/monitor/internal/signer"
	"github.com/artem-benda/monitor/internal/version"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	store        storage.Storage
	dbpool       *pgxpool.Pool
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	version.PrintVersion(buildVersion, buildDate, buildCommit)
	parseFlags()
	var err error

	if err = logger.Initialize(config.LogLevel); err != nil {
		panic(err)
	} else {
		defer func() {
			err = logger.Log.Sync()
			if err != nil {
				logger.Log.Error("Could not sync log data", zap.Error(err))
			}
		}()
	}

	var flushStorage func() error

	if config.DatabaseDSN != "" {
		dbpool = newConnectionPool(config.DatabaseDSN)
		initDB(dbpool)
		defer dbpool.Close()

		retryController := retry.NewRetryController(errors.ErrStorageConnection{})
		store = storage.NewDBStorage(dbpool, retryController)
	} else {
		store, flushStorage, err = storage.NewMemStorage(config.StoreIntervalSeconds, config.StoreFileName, config.StoreRestoreFromFile)
		if err != nil {
			panic(err)
		}
	}

	r := newAppRouter()
	// Настройки сервера
	server := &http.Server{Addr: config.Endpoint, Handler: r}

	// Контекст сервера
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Слушаем нужные сигналы от ОС
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sig

		// Даем максимум 30 секунд на выполнение останова
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			defer shutdownCancel()
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Выполняем graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Запускаем сервер
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Ожидаем завершения
	<-serverCtx.Done()
	// Сбрасываем на диск данные из хранилища, только для memStorage
	if flushStorage != nil {
		err = flushStorage()
		logger.Log.Error("error flushing storage on shutdown")
	}
}

func newAppRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.LoggerMiddleware)
	r.Use(ipfilter.NewIPFilterMiddleware(mustParseTrustedSubnetCIDR(config.TrustedSubnet)))
	r.Use(signer.CreateVerifyAndSignMiddleware([]byte(config.Key)))
	r.Use(crypt.NewDecryptMiddleware(mustParseRSAPrivateKey(config.RSAPrivKeyBase64)))
	r.Use(gzipper.GzipMiddleware)

	r.Mount("/debug", middleware.Profiler())

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.MakeGetAllHandler(store))
		r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.MakeUpdatePathHandler(store))
		r.Post("/update/", handlers.MakeUpdateJSONHandler(store))
		r.Post("/updates/", handlers.MakeUpdateBatchJSONHandler(store))
		r.Get("/ping", handlers.MakePingDatabaseHandler(dbpool))
		r.Get("/value/{metricType}/{metricName}", handlers.MakeGetHandler(store))
		r.Post("/value/", handlers.MakeGetJSONHandler(store))
	})
	return r
}

func initDB(dbpool *pgxpool.Pool) {
	createTblMetrics := `CREATE TABLE IF NOT EXISTS metrics(
		mtype text NOT NULL,
		mname text NOT NULL,
		gauge double precision,
		counter bigint,
		PRIMARY KEY (mtype, mname)
		)`
	_, err := dbpool.Exec(context.Background(), createTblMetrics)
	if err != nil {
		panic(err)
	}
}
