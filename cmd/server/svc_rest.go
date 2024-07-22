package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artem-benda/monitor/internal/crypt"
	"github.com/artem-benda/monitor/internal/gzipper"
	"github.com/artem-benda/monitor/internal/ipfilter"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/server/handlers"
	"github.com/artem-benda/monitor/internal/signer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func mustRunRestServer(flushStorage func() error) {
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
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Ожидаем завершения
	<-serverCtx.Done()
	// Сбрасываем на диск данные из хранилища, только для memStorage
	if flushStorage != nil {
		err = flushStorage()
		if err != nil {
			logger.Log.Error("error flushing storage on shutdown")
		}
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
