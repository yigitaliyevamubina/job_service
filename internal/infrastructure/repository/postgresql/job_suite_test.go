package postgresql

import (
	"context"
	"fifth_exam/job_service/internal/entity"
	"fifth_exam/job_service/internal/pkg/config"
	"fifth_exam/job_service/internal/pkg/postgres"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type JobTestSite struct {
	suite.Suite
	Repository  *JobRepo
	CleanUpFunc func()
}

func (s *JobTestSite) SetupSuite() {
	pgPool, _ := postgres.New(config.New())
	s.Repository = NewJobRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}

func (j *JobTestSite) TestJobCRUD() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

	job := &entity.Job{
		Id:          uuid.New().String(),
		Title:       "New title",
		Description: "New Description",
		OwnerId:     "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		Price:       12412.12,
		FromDate:    "2023-12-12",
		ToDate:      "2023-12-27",
	}

	// Create
	createReq, err := j.Repository.Create(ctx, job)
	j.Suite.NoError(err)
	j.Suite.NotNil(createReq)

	// Get
	getRes, err := j.Repository.Get(ctx, "id", job.Id)
	j.Suite.NoError(err)
	j.Suite.NotNil(getRes)
	j.Suite.Equal(job.Title, getRes.Title)
	j.Suite.Equal(job.Description, getRes.Description)
	j.Suite.Equal(job.OwnerId, getRes.OwnerId)
	j.Suite.Equal(job.Price, getRes.Price)
	j.Suite.Equal(job.FromDate, getRes.FromDate)
	j.Suite.Equal(job.ToDate, getRes.ToDate)

	// List
	listRes, err := j.Repository.List(ctx, &entity.GetListFilter{Page: 1, Limit: 10})
	j.Suite.NoError(err)
	j.Suite.NotNil(listRes)

	// Update
	job.Description = "Updated Description"
	err = j.Repository.Update(ctx, job)
	j.Suite.NoError(err)

	// Verify updated description
	getRes, err = j.Repository.Get(ctx, "id", job.Id)
	j.Suite.NoError(err)
	j.Suite.NotNil(getRes)
	j.Suite.Equal("Updated Description", getRes.Description)

	// Delete
	err = j.Repository.Delete(ctx, "id", job.Id)
	j.Suite.NoError(err)
}

func (s *JobTestSite) TearDownSuite() {
	s.CleanUpFunc()
}

func TestDoctorNotesTestSuite(t *testing.T) {
	suite.Run(t, new(JobTestSite))
}
