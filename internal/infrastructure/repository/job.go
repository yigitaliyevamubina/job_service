package repository

import (
	"context"
	"fifth_exam/job_service/internal/entity"
)

type Job interface {
	Create(ctx context.Context, req *entity.Job) (*entity.Job, error)
	Get(ctx context.Context, field, value string) (*entity.Job, error)
	List(ctx context.Context, req *entity.GetListFilter) ([]*entity.Job, error)
	Update(ctx context.Context, req *entity.Job) error
	Delete(ctx context.Context, field, value string) error
}
