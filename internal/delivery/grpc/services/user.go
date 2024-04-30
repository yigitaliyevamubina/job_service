package services

import (
	"context"
	pb "fifth_exam/job_service/genproto/job_service"
	"fifth_exam/job_service/internal/entity"
	"fifth_exam/job_service/internal/pkg/otlp"
	"fifth_exam/job_service/internal/usecase"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	serviceNameJob = "jobService"
	spanNameJob    = "jobUsecase"
)

type jobRPC struct {
	logger     *zap.Logger
	jobUsecase usecase.Job
}

func NewRPC(logger *zap.Logger, jobUsecase usecase.Job) pb.JobServiceServer {
	return &jobRPC{
		logger:     logger,
		jobUsecase: jobUsecase,
	}
}

func (j *jobRPC) Create(ctx context.Context, in *pb.Job) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"Create")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "delivery", Value: attribute.StringValue("CreatingJob")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	id := uuid.New().String()
	_, err := j.jobUsecase.Create(ctxWithSpan, &entity.Job{
		Id:          id,
		Title:       in.Title,
		OwnerId:     in.OwnerId,
		Price:       in.Price,
		Description: in.Description,
		CreatedAt:   time.Now(),
		FromDate:    in.FromDate,
		ToDate:      in.ToDate,
	})
	if err != nil {
		j.logger.Error("jobUseCase.Create", zap.Error(err))
		return &pb.Job{}, status.Errorf(codes.Internal, "failed to create job: %v", err)
	}
	in.Id = id
	return in, nil
}

func (j *jobRPC) Update(ctx context.Context, in *pb.Job) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"Update")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "delivery", Value: attribute.StringValue("UpdatingJob")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	err := j.jobUsecase.Update(ctxWithSpan, &entity.Job{
		Id:          in.Id,
		Title:       in.Title,
		OwnerId:     in.OwnerId,
		Price:       in.Price,
		Description: in.Description,
		CreatedAt:   time.Now(),
		FromDate:    in.FromDate,
		ToDate:      in.ToDate,
	})
	if err != nil {
		j.logger.Error("jobUseCase.Update", zap.Error(err))
		return &pb.Job{}, status.Errorf(codes.Internal, "failed to update job: %v", err)
	}

	return in, nil
}

func (j *jobRPC) Get(ctx context.Context, in *pb.JobRequest) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"Get")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "delivery", Value: attribute.StringValue("GettingJob")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	job, err := j.jobUsecase.Get(ctxWithSpan, in.Field, in.Value)
	if err != nil {
		j.logger.Error("jobUseCase.Get", zap.Error(err))
		return &pb.Job{}, status.Errorf(codes.Internal, "failed to get job: %v", err)
	}

	return &pb.Job{
		Id:          job.Id,
		Title:       job.Title,
		OwnerId:     job.OwnerId,
		Price:       job.Price,
		Description: job.Description,
		CreatedAt:   job.CreatedAt.String(),
		UpdatedAt:   job.UpdatedAt.String(),
		FromDate:    job.FromDate,
		ToDate:      job.ToDate,
	}, nil
}

func (j *jobRPC) Delete(ctx context.Context, in *pb.JobRequest) (*empty.Empty, error) {
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"Delete")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "delivery", Value: attribute.StringValue("DeletingJob")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	err := j.jobUsecase.Delete(ctxWithSpan, in.Field, in.Value)
	if err != nil {
		j.logger.Error("jobUseCase.Delete", zap.Error(err))
		return &empty.Empty{}, status.Errorf(codes.Internal, "failed to delete job: %v", err)
	}

	return &empty.Empty{}, nil
}

func (j *jobRPC) GetList(ctx context.Context, in *pb.GetListFilter) (*pb.Jobs, error) {
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"List")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "delivery", Value: attribute.StringValue("GettingJobList")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	filter := &entity.GetListFilter{
		Limit:   in.Limit,
		Page:    in.Page,
		OrderBy: in.OrderBy,
	}

	jobs, err := j.jobUsecase.List(ctxWithSpan, filter)
	if err != nil {
		j.logger.Error("jobUseCase.List", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to retrieve jobs: %v", err)
	}

	var pbJobs pb.Jobs
	for _, job := range jobs {
		pbJob := &pb.Job{
			Id:          job.Id,
			Title:       job.Title,
			OwnerId:     job.OwnerId,
			Price:       job.Price,
			Description: job.Description,
			CreatedAt:   job.CreatedAt.String(),
			UpdatedAt:   job.UpdatedAt.String(),
			FromDate:    job.FromDate,
			ToDate:      job.ToDate,
		}
		pbJobs.Jobs = append(pbJobs.Jobs, pbJob)
		pbJobs.Count++
	}

	return &pbJobs, nil
}
