package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	"github.com/nickolation/grpc-task-hezzl/configs"
	"github.com/nickolation/grpc-task-hezzl/storage/db"
)

func main() {
	errChan := make(chan error, 0)
	go func()  {
		go func() {
			logrus.Fatalf("error with the setupping and setting the grpc-server - [%v]", <-errChan)
		}()
	}()
	if err := godotenv.Load(); err != nil {
		errChan <- err
	}

	cfg, err := configs.NewGrpcConfigManager()
	if err != nil {
		errChan <- err
	}

	d, err := db.ConnectToPostgresDb(
		db.WithPort(cfg.GetDbPort()),
		db.WithHost(cfg.GetDbHost()),
		db.WithDbName(cfg.GetDbName()),
		db.WithDisabledSSLMode(),
		db.WithPassword(os.Getenv("POSTGRES_PASSWORD")),
		db.WithUsername(cfg.GetDbUsername()),
	)

	if err != nil {
		errChan <- err
	}

	db.UpdatePostgresDbScheme(d)
}
