package models

import (
	"context"
	"testing"
	"time"

	"github.com/m1yon/jobsummoner/internal/database"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestJobs(t *testing.T) {
	jobsToCreate := []Job{
		{
			Position:      "Software Developer",
			URL:           "https://linkedin.com/jobs/1",
			Location:      "San Francisco",
			SourceID:      "linkedin",
			CompanyID:     "/google",
			CompanyName:   "Google",
			CompanyAvatar: "https://upload.wikimedia.org/wikipedia/commons/c/c1/Google_%22G%22_logo.svg",
		},
		{
			Position:      "Manager",
			URL:           "https://linkedin.com/jobs/2",
			Location:      "Seattle",
			SourceID:      "linkedin",
			CompanyID:     "/microsoft",
			CompanyName:   "Microsoft",
			CompanyAvatar: "https://blogs.microsoft.com/wp-content/uploads/prod/2012/08/8867.Microsoft_5F00_Logo_2D00_for_2D00_screen.jpg",
		},
	}

	t.Run("create job and immediately get created job", func(t *testing.T) {
		ctx := context.Background()
		companies, jobs := newTestModels(t)

		jobToCreate := jobsToCreate[0]

		id, err := jobs.Create(ctx, jobToCreate)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		assertCompanyExist(t, companies, jobToCreate.CompanyID)

		job, err := jobs.Get(ctx, id)
		assert.NoError(t, err)
		assertJobsAreEqual(t, jobToCreate, job)
	})

	t.Run("CreateJobs returns a list of job IDs", func(t *testing.T) {
		ctx := context.Background()
		_, jobs := newTestModels(t)

		ids, errs := jobs.CreateMany(ctx, jobsToCreate)
		assert.Equal(t, 0, len(errs))

		for i := range jobsToCreate {
			id := ids[i]
			assert.NotEmpty(t, id)
		}
	})

	t.Run("new companies exist after jobs created", func(t *testing.T) {
		ctx := context.Background()
		companies, jobs := newTestModels(t)

		ids, errs := jobs.CreateMany(ctx, jobsToCreate)
		assert.Equal(t, len(jobsToCreate), len(ids))
		assert.Equal(t, 0, len(errs))

		for _, jobToCreate := range jobsToCreate {
			assertCompanyExist(t, companies, jobToCreate.CompanyID)
		}
	})

	t.Run("can query new companies after jobs created", func(t *testing.T) {
		ctx := context.Background()
		companies, jobs := newTestModels(t)

		ids, errs := jobs.CreateMany(ctx, jobsToCreate)
		assert.Equal(t, len(jobsToCreate), len(ids))
		assert.Equal(t, 0, len(errs))

		for _, jobToCreate := range jobsToCreate {
			createdCompany, err := companies.Get(ctx, jobToCreate.CompanyID)
			assert.NoError(t, err)

			assertCompanyFieldsOnJob(t, jobToCreate, createdCompany)
		}
	})

	t.Run("can get jobs after jobs created", func(t *testing.T) {
		ctx := context.Background()
		_, jobs := newTestModels(t)

		ids, errs := jobs.CreateMany(ctx, jobsToCreate)
		assert.Equal(t, len(jobsToCreate), len(ids))
		assert.Equal(t, 0, len(errs))

		for i, jobToCreate := range jobsToCreate {
			id := ids[i]
			job, err := jobs.Get(ctx, id)
			assert.NoError(t, err)

			assert.EqualExportedValues(t, jobToCreate, job)
		}
	})

	t.Run("get jobs", func(t *testing.T) {
		ctx := context.Background()
		_, jobs := newTestModels(t)

		jobs.CreateMany(ctx, jobsToCreate)

		res, errs := jobs.GetMany(ctx)
		assert.Empty(t, errs)
		assert.Equal(t, jobsToCreate, res)
	})
}

func assertCompanyFieldsOnJob(t *testing.T, job Job, company Company) {
	assert.Equal(t, company.ID, job.CompanyID)
	assert.Equal(t, company.Name, job.CompanyName)
	assert.Equal(t, company.Name, job.CompanyName)
	assert.Equal(t, company.Avatar, job.CompanyAvatar)
	assert.Equal(t, company.SourceID, job.SourceID)
}

func assertJobsAreEqual(t *testing.T, expectedJob, actualJob Job) {
	assert.Equal(t, expectedJob.Position, actualJob.Position)
	assert.Equal(t, expectedJob.URL, actualJob.URL)
	assert.Equal(t, expectedJob.Location, actualJob.Location)
	assert.Equal(t, expectedJob.SourceID, actualJob.SourceID)

	assert.Equal(t, expectedJob.CompanyID, actualJob.CompanyID)
	assert.Equal(t, expectedJob.CompanyAvatar, actualJob.CompanyAvatar)
	assert.Equal(t, expectedJob.CompanyName, actualJob.CompanyName)
	assert.Equal(t, expectedJob.CompanyURL, actualJob.CompanyURL)

	assert.Equal(t, expectedJob.LastPosted.Round(time.Second).UTC(), actualJob.LastPosted.Round(time.Second).UTC())
}

func newTestModels(t *testing.T) (*CompanyModel, *JobModel) {
	db, err := database.NewInMemoryDB()

	if err != nil {
		t.Fatal(err)
	}

	queries := database.New(db)
	companies := &CompanyModel{queries}
	jobs := &JobModel{Queries: queries, Companies: companies}

	return companies, jobs
}
