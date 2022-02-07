package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nickolation/grpc-task-hezzl/configs"
	api "github.com/nickolation/grpc-task-hezzl/grpc/grpc_client_api"
	"github.com/sirupsen/logrus"
)

func main() {
	errChan := make(chan error, 0)
	go func() {
		logrus.Fatalf("error with the setupping and setting the api-client - [%v]", <-errChan)
	}()

	if err := godotenv.Load(); err != nil {
		errChan <- err
	}

	cfg, err := configs.NewGrpcConfigManager()
	if err != nil {
		errChan <- err
	}

	clientServer := api.NewClientAPIServer(cfg)
	go func() {
		if err := clientServer.RunClientServer(); err != nil {
			logrus.Errorf("error with the running the http server %s", err.Error())
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit 

	if err := clientServer.ShutdownClientServer(context.Background()); err != nil {
		logrus.Fatalf("error with the shotdowning the http server %s", err.Error())
	}
}
