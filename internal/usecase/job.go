package usecase

import (
	"context"
	"fifth_exam/job_service/internal/entity"
	"fifth_exam/job_service/internal/infrastructure/repository"
	"fifth_exam/job_service/internal/pkg/otlp"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	serviceNameJob = "jobService"
	spanNameJob    = "jobUsecase"
)

type Job interface {
	Create(ctx context.Context, req *entity.Job) (*entity.Job, error)
	Get(ctx context.Context, field, value string) (*entity.Job, error)
	List(ctx context.Context, req *entity.GetListFilter) ([]*entity.Job, error)
	Update(ctx context.Context, req *entity.Job) error
	Delete(ctx context.Context, field, value string) error
}

type jobService struct {
	BaseUseCase
	repo       repository.Job
	ctxTimeout time.Duration
}

func NewJobService(ctxTimeout time.Duration, repo repository.Job) Job {
	return &jobService{
		repo:       repo,
		ctxTimeout: ctxTimeout,
	}
}

func (j *jobService) Create(ctx context.Context, req *entity.Job) (*entity.Job, error) {
	ctx, cancel := context.WithTimeout(ctx, j.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"Create")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "repo", Value: attribute.StringValue("Creating job")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	return j.repo.Create(ctxWithSpan, req)
}

func (j *jobService) Get(ctx context.Context, field, value string) (*entity.Job, error) {
	ctx, cancel := context.WithTimeout(ctx, j.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"Get")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "usecase", Value: attribute.StringValue("Getting job")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	return j.repo.Get(ctxWithSpan, field, value)
}

func (j *jobService) List(ctx context.Context, req *entity.GetListFilter) ([]*entity.Job, error) {
	ctx, cancel := context.WithTimeout(ctx, j.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"List")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "usecase", Value: attribute.StringValue("Getting job list")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	return j.repo.List(ctxWithSpan, req)
}

func (j *jobService) Update(ctx context.Context, req *entity.Job) error {
	ctx, cancel := context.WithTimeout(ctx, j.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"Update")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "usecase", Value: attribute.StringValue("Updating job")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	return j.repo.Update(ctxWithSpan, req)
}

func (j *jobService) Delete(ctx context.Context, field, value string) error {
	ctx, cancel := context.WithTimeout(ctx, j.ctxTimeout)
	defer cancel()
	ctx, span := otlp.Start(ctx, serviceNameJob, spanNameJob+"Delete")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "usecase", Value: attribute.StringValue("Deleting job")})
	ctxWithSpan := trace.ContextWithSpan(ctx, span)

	return j.repo.Delete(ctxWithSpan, field, value)
}
