package postgresql

import (
	"context"
	"database/sql"
	"fifth_exam/job_service/internal/entity"
	"fifth_exam/job_service/internal/pkg/otlp"
	"fifth_exam/job_service/internal/pkg/postgres"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"go.opentelemetry.io/otel/attribute"
)

const (
	jobsTableName     = "jobs"
	jobServiceName    = "jobService"
	jobSpanRepoPrefix = "jobServiceRepo"
)

type JobRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewJobRepo(db *postgres.PostgresDB) *JobRepo {
	return &JobRepo{
		tableName: jobsTableName,
		db:        db,
	}
}

func (j *JobRepo) Create(ctx context.Context, req *entity.Job) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"Create")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "sql", Value: attribute.StringValue("Creating job")})
	span.AddEvent("starting creating job")

	data := map[string]interface{}{
		"id":          req.Id,
		"title":       req.Title,
		"description": req.Description,
		"owner_id":    req.OwnerId,
		"price":       req.Price,
		"created_at":  req.CreatedAt,
		"from_date":   req.FromDate,
		"to_date":     req.ToDate,
	}

	query, args, err := j.db.Sq.Builder.Insert(j.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, j.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", j.tableName, "create"))
	}

	_, err = j.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, j.db.Error(err)
	}

	return req, nil
}

func (j *JobRepo) Get(ctx context.Context, field, value string) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"Get")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "sql", Value: attribute.StringValue("Getting job")})

	var job entity.Job

	queryBuilder := j.jobsSelectQueryPrefix().Where(
		squirrel.And{
			squirrel.Eq{field: value},
			squirrel.Eq{"deleted_at": nil},
		},
	)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Object not found")
		return nil, j.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", j.tableName, "get"))
	}

	var (
		deletedAt sql.NullTime
		updatedAt sql.NullTime
	)
	if err = j.db.QueryRow(ctx, query, args...).Scan(
		&job.Id,
		&job.Title,
		&job.Description,
		&job.OwnerId,
		&job.Price,
		&job.CreatedAt,
		&updatedAt,
		&deletedAt,
		&job.FromDate,
		&job.ToDate,
	); err != nil {
		return nil, j.db.Error(err)
	}

	if updatedAt.Valid {
		job.UpdatedAt = updatedAt.Time
	}

	return &job, nil
}

func (j *JobRepo) List(ctx context.Context, req *entity.GetListFilter) ([]*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"List")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "sql", Value: attribute.StringValue("Getting job list")})

	var jobs []*entity.Job

	queryBuilder := j.jobsSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	if req.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))
	}

	if req.OrderBy != "" {
		queryBuilder = queryBuilder.OrderBy(req.OrderBy)
	}

	if !req.IncludeDeleted {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"deleted_at": nil})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, j.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", j.tableName, "list"))
	}

	rows, err := j.db.Query(ctx, query, args...)
	if err != nil {
		return nil, j.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var job entity.Job
		var deletedAt sql.NullTime
		var updatedAt sql.NullTime
		if err = rows.Scan(
			&job.Id,
			&job.Title,
			&job.Description,
			&job.OwnerId,
			&job.Price,
			&job.CreatedAt,
			&updatedAt,
			&deletedAt,
			&job.FromDate,
			&job.ToDate,
		); err != nil {
			return nil, j.db.Error(err)
		}

		if updatedAt.Valid {
			job.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			job.DeletedAt = deletedAt.Time
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (j *JobRepo) Update(ctx context.Context, req *entity.Job) error {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"Update")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "sql", Value: attribute.StringValue("Updating job")})

	data := map[string]interface{}{
		"title":       req.Title,
		"description": req.Description,
		"owner_id":    req.OwnerId,
		"price":       req.Price,
		"updated_at":  req.UpdatedAt,
		"from_date":   req.FromDate,
		"to_date":     req.ToDate,
	}

	sqlStr, args, err := j.db.Sq.Builder.
		Update(j.tableName).
		SetMap(data).
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
	if err != nil {
		return j.db.ErrSQLBuild(err, j.tableName+" update")
	}

	commandTag, err := j.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return j.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return j.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (j *JobRepo) Delete(ctx context.Context, field, value string) error {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"SoftDelete")
	defer span.End()

	span.SetAttributes(attribute.KeyValue{Key: "sql", Value: attribute.StringValue("Soft deleting job")})

	sqlStr, args, err := j.db.Sq.Builder.
		Update(j.tableName).
		Set("deleted_at", time.Now()).
		Where(j.db.Sq.Equal(field, value)).
		ToSql()
	if err != nil {
		return j.db.ErrSQLBuild(err, j.tableName+" soft delete")
	}

	_, err = j.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return j.db.Error(err)
	}

	return nil
}

func (j *JobRepo) jobsSelectQueryPrefix() squirrel.SelectBuilder {
	return j.db.Sq.Builder.Select(
		"id",
		"title",
		"description",
		"owner_id",
		"price",
		"created_at",
		"updated_at",
		"deleted_at",
		"from_date",
		"to_date",
	).From(j.tableName)
}
