package app

import (
	"fifth_exam/job_service/internal/delivery/kafka/handlers"
	"fifth_exam/job_service/internal/infrastructure/kafka"
	"fifth_exam/job_service/internal/infrastructure/repository/postgresql"
	"fifth_exam/job_service/internal/pkg/config"
	"fifth_exam/job_service/internal/pkg/postgres"
	"fifth_exam/job_service/internal/usecase"
	"fifth_exam/job_service/internal/usecase/event"
	"fmt"
	"time"

	logpkg "fifth_exam/job_service/internal/pkg/logger"

	"go.uber.org/zap"
)

type JobConsumer struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	BrokerConsumer event.BrokerConsumer
}

func NewJobConsumer(conf *config.Config) (*JobConsumer, error) {
	logger, err := logpkg.New(conf.LogLevel, conf.Environment, conf.APP+"_cli"+".lo")
	if err != nil {
		return nil, err
	}

	consumer := kafka.NewConsumer(logger)

	db, err := postgres.New(conf)
	if err != nil {
		return nil, err
	}

	return &JobConsumer{Config: conf, Logger: logger, DB: db, BrokerConsumer: consumer}, nil
}

func (u *JobConsumer) Run() error {

	// repo init
	jobRepo := postgresql.NewJobRepo(u.DB)

	// usecase init
	duration, err := time.ParseDuration(u.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}
	jobUseCase := usecase.NewJobService(duration, jobRepo)

	// event handler
	eventHandler := handlers.NewUserConsumerHandler(u.Config, u.BrokerConsumer, u.Logger, jobUseCase)

	return eventHandler.HandlerEvents()
}

func (u *JobConsumer) Close() {
	u.BrokerConsumer.Close()

	u.Logger.Sync()
}
