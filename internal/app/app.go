package app

import (
	pb "fifth_exam/job_service/genproto/job_service"
	"fifth_exam/job_service/internal/delivery/grpc/server"
	"fifth_exam/job_service/internal/delivery/grpc/services"
	grpc_service_clients "fifth_exam/job_service/internal/infrastructure/grpc_service_client"
	"fifth_exam/job_service/internal/infrastructure/kafka"
	"fifth_exam/job_service/internal/infrastructure/repository/postgresql"
	"fifth_exam/job_service/internal/pkg/config"
	"fifth_exam/job_service/internal/pkg/logger"
	"fifth_exam/job_service/internal/pkg/otlp"
	"fifth_exam/job_service/internal/pkg/postgres"
	"fifth_exam/job_service/internal/usecase"
	"fifth_exam/job_service/internal/usecase/event"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	ServiceClients grpc_service_clients.ServiceClients
	GrpcServer     *grpc.Server
	ShutdownOTLP   func() error
	BrokerConsumer event.BrokerConsumer
}

func NewApp(cfg *config.Config) (*App, error) {
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	shutdownOTLP, err := otlp.InitOTLPProvider(cfg)
	if err != nil {
		fmt.Println("initOTLP provider error")
		return nil, err
	}
	grpcServer := grpc.NewServer()
	clients, err := grpc_service_clients.New(cfg)
	if err != nil {
		return nil, err
	}
	brokerConsumer := kafka.NewConsumer(logger)

	return &App{
		Config:         cfg,
		Logger:         logger,
		DB:             db,
		GrpcServer:     grpcServer,
		ServiceClients: clients,
		BrokerConsumer: brokerConsumer,
		ShutdownOTLP:   shutdownOTLP,
	}, nil
}

func (a *App) Run() error {
	var (
		contextTimeout time.Duration
	)
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}

	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	jobRepo := postgresql.NewJobRepo(a.DB)

	jobUseCase := usecase.NewJobService(contextTimeout, jobRepo)

	pb.RegisterJobServiceServer(a.GrpcServer, services.NewRPC(a.Logger, jobUseCase))

	// a.BrokerConsumer.Run()

	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))
	if err := server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}
	return nil
}

func (a *App) Stop() {
	// closing client service connections
	a.ServiceClients.Close()
	// stop gRPC server
	a.GrpcServer.Stop()

	// database connection
	a.DB.Close()

	// broker consumer connection
	a.BrokerConsumer.Close()

	// zap logger sync
	a.Logger.Sync()
}
