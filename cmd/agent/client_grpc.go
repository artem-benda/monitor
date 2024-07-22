package main

import (
	"log"

	pb "github.com/artem-benda/monitor/internal/grpc/mon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func mustCreateGRPCClient() (pb.MonitorServiceClient, *grpc.ClientConn) {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(config.ServerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	// получаем переменную интерфейсного типа MonitorServiceClient,
	// через которую будем отправлять сообщения
	c := pb.NewMonitorServiceClient(conn)
	return c, conn
}
