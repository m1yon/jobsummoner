package job

import (
	"context"
	"testing"
	"time"

	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestSqliteJobService(t *testing.T) {
	jobsToCreate := []jobsummoner.Job{
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
		db, _ := sqlitedb.NewInMemoryDB()

		queries := sqlitedb.New(db)
		companyService := company.NewDefaultCompanyService(queries)
		jobService := NewDefaultJobService(queries, companyService)

		jobToCreate := jobsummoner.Job{
			Position:    "Software Developer",
			URL:         "https://linkedin.com/jobs/1",
			Location:    "San Francisco",
			CompanyID:   "/google",
			CompanyName: "Google",
			SourceID:    "linkedin",
		}

		id, err := jobService.CreateJob(ctx, jobToCreate)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		doesCompanyExist, err := companyService.DoesCompanyExist(ctx, jobToCreate.CompanyID)
		assert.NoError(t, err)
		assert.Equal(t, true, doesCompanyExist)

		job, err := jobService.GetJob(ctx, id)
		assert.NoError(t, err)
		assertJobsAreEqual(t, jobToCreate, job)
	})

	t.Run("CreateJobs returns a list of job IDs", func(t *testing.T) {
		ctx := context.Background()
		db, _ := sqlitedb.NewInMemoryDB()

		queries := sqlitedb.New(db)
		companyService := company.NewDefaultCompanyService(queries)
		jobService := NewDefaultJobService(queries, companyService)

		ids, errs := jobService.CreateJobs(ctx, jobsToCreate)
		assert.Equal(t, 0, len(errs))

		for i := range jobsToCreate {
			id := ids[i]
			assert.NotEmpty(t, id)
		}
	})

	t.Run("new companies exist after jobs created", func(t *testing.T) {
		ctx := context.Background()
		db, _ := sqlitedb.NewInMemoryDB()

		queries := sqlitedb.New(db)
		companyService := company.NewDefaultCompanyService(queries)
		jobService := NewDefaultJobService(queries, companyService)

		ids, errs := jobService.CreateJobs(ctx, jobsToCreate)
		assert.Equal(t, len(jobsToCreate), len(ids))
		assert.Equal(t, 0, len(errs))

		for _, jobToCreate := range jobsToCreate {
			doesCompanyExist, err := companyService.DoesCompanyExist(ctx, jobToCreate.CompanyID)
			assert.NoError(t, err)
			assert.Equal(t, true, doesCompanyExist)
		}
	})

	t.Run("can query new companies after jobs created", func(t *testing.T) {
		ctx := context.Background()
		db, _ := sqlitedb.NewInMemoryDB()

		queries := sqlitedb.New(db)
		companyService := company.NewDefaultCompanyService(queries)
		jobService := NewDefaultJobService(queries, companyService)

		ids, errs := jobService.CreateJobs(ctx, jobsToCreate)
		assert.Equal(t, len(jobsToCreate), len(ids))
		assert.Equal(t, 0, len(errs))

		for _, jobToCreate := range jobsToCreate {
			createdCompany, err := companyService.GetCompany(ctx, jobToCreate.CompanyID)
			assert.NoError(t, err)

			assert.Equal(t, jobToCreate.CompanyID, createdCompany.ID)
			assert.Equal(t, jobToCreate.CompanyName, createdCompany.Name)
			assert.Equal(t, jobToCreate.CompanyName, createdCompany.Name)
			assert.Equal(t, jobToCreate.CompanyAvatar, createdCompany.Avatar)
			assert.Equal(t, jobToCreate.SourceID, createdCompany.SourceID)
		}
	})

	t.Run("can get jobs after jobs created", func(t *testing.T) {
		ctx := context.Background()
		db, _ := sqlitedb.NewInMemoryDB()

		queries := sqlitedb.New(db)
		companyService := company.NewDefaultCompanyService(queries)
		jobService := NewDefaultJobService(queries, companyService)

		ids, errs := jobService.CreateJobs(ctx, jobsToCreate)
		assert.Equal(t, len(jobsToCreate), len(ids))
		assert.Equal(t, 0, len(errs))

		for i, jobToCreate := range jobsToCreate {
			id := ids[i]
			job, err := jobService.GetJob(ctx, id)
			assert.NoError(t, err)

			assert.Equal(t, jobToCreate.Position, job.Position)
			assert.Equal(t, jobToCreate.Location, job.Location)
			assert.Equal(t, jobToCreate.SourceID, job.SourceID)
			assert.Equal(t, jobToCreate.URL, job.URL)
		}
	})

	t.Run("get jobs", func(t *testing.T) {
		ctx := context.Background()
		db, _ := sqlitedb.NewInMemoryDB()

		queries := sqlitedb.New(db)
		companyService := company.NewDefaultCompanyService(queries)
		jobService := NewDefaultJobService(queries, companyService)

		jobService.CreateJobs(ctx, jobsToCreate)

		jobs, errs := jobService.GetJobs(ctx)
		assert.Empty(t, errs)
		assert.Equal(t, jobsToCreate, jobs)
	})
}

func assertJobsAreEqual(t *testing.T, expectedJob, actualJob jobsummoner.Job) {
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
