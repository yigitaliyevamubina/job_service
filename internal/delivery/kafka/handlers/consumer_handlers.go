package handlers

import (
	"context"
	"encoding/json"
	"fifth_exam/job_service/internal/entity"
	"fifth_exam/job_service/internal/infrastructure/kafka"
	"fifth_exam/job_service/internal/pkg/config"
	"fifth_exam/job_service/internal/usecase"
	"fifth_exam/job_service/internal/usecase/event"
	pb 	"fifth_exam/job_service/genproto/job_service"
	"fmt"
	"time"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

type jobConsumerHandler struct {
	config         *config.Config
	brokerConsumer event.BrokerConsumer
	logger         *zap.Logger
	jobUsecase     usecase.Job
}

func NewUserConsumerHandler(conf *config.Config,
	brokerConsumer event.BrokerConsumer,
	logger *zap.Logger,
	jobUsecase usecase.Job) *jobConsumerHandler {
	return &jobConsumerHandler{
		config:         conf,
		brokerConsumer: brokerConsumer,
		logger:         logger,
		jobUsecase:     jobUsecase,
	}
}

func (u *jobConsumerHandler) HandlerEvents() error {
	consumerConfig := kafka.NewConsumerConfig(
		u.config.Kafka.Address,
		u.config.Kafka.Topic.JobTopic,
		"1",
		func(ctx context.Context, key, value []byte) error {
			var job pb.Job

			if err := json.Unmarshal(value, &job); err != nil {
				return err
			}

			pp.Println(job)

			req := entity.Job{
				Id:          job.Id,
				Title:       job.Title,
				Description: job.Description,
				OwnerId:     job.OwnerId,
				Price:       job.Price,
			}

			ctxNew, err := context.WithTimeout(context.Background(), time.Second*7)
			if err != nil {
				fmt.Println(err, "Context error")
			}
			_, errr := u.jobUsecase.Create(ctxNew, &req)
			if errr != nil {
				fmt.Println(errr, "Create Job By Kafka Error")
			}

			return nil
		},
	)

	u.brokerConsumer.RegisterConsumer(consumerConfig)
	u.brokerConsumer.Run()

	return nil
}
