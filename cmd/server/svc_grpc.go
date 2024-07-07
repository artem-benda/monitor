package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/artem-benda/monitor/internal/grpc/mon"
	"github.com/artem-benda/monitor/internal/logger"
	"github.com/artem-benda/monitor/internal/server/grpc"
	"go.uber.org/zap"
	g "google.golang.org/grpc"
)

func mustRunGrpcServer(flushStorage func() error) {
	listen, err := net.Listen("tcp", config.Endpoint)
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := g.NewServer()
	// регистрируем сервис
	pb.RegisterMonitorServiceServer(s, &grpc.MetricsGrpsServer{})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {

		logger.Log.Debug("server sut down", zap.Error(err))
		// Сбрасываем на диск данные из хранилища, только для memStorage
		if flushStorage != nil {
			err = flushStorage()
			if err != nil {
				logger.Log.Error("error flushing storage on shutdown")
			}
		}
	}
}
