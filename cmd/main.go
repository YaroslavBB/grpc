package main

import (
	"grpc/app/config"
	"grpc/app/external"
	"grpc/app/internal/repository"
	"grpc/app/internal/usecase"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	log := NewNoFileLogger("grpc")

	config, err := config.NewConfig("../config/config.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	db, err := sqlx.Open("postgres", config.DbConn())
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	imageRepo := repository.NewImageRepository(db)
	imageUsecase := usecase.NewImageUsecase(imageRepo, log)
	grpcServer := external.NewGRPCServer(imageUsecase, log, config)

	grpcServer.Run()

}

func textFormatter() *logrus.TextFormatter {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "02.01.2006 15:04:05"
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true

	return customFormatter
}

// NewNoFileLogger файловый логгер без реальной записи в файл
func NewNoFileLogger(module string) *logrus.Logger {
	newLogger := logrus.New()
	newLogger.Formatter = textFormatter()

	newLogger.Level = logrus.DebugLevel

	return newLogger
}
