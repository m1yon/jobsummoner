package tests

import (
	"context"
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/stretchr/testify/assert"
)

func TestSqliteJobService(t *testing.T) {
	t.Run("create job and immediately get created job", func(t *testing.T) {
		ctx := context.Background()
		db := sqlitedb.NewTestDB()
		companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
		companyService := company.NewDefaultCompanyService(companyRepository)
		jobRepository := sqlitedb.NewSqliteJobRepository(db)
		jobService := job.NewDefaultJobService(jobRepository, companyService)

		jobToCreate := jobsummoner.Job{
			Position:    "Software Developer",
			URL:         "https://linkedin.com/jobs/1",
			Location:    "San Francisco",
			CompanyID:   "/google",
			CompanyName: "Google",
			SourceID:    "1",
		}

		id, err := jobService.CreateJob(ctx, jobToCreate)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		doesCompanyExist, err := companyService.DoesCompanyExist(ctx, jobToCreate.CompanyID)
		assert.NoError(t, err)
		assert.Equal(t, true, doesCompanyExist)

		job, err := jobService.GetJob(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, jobToCreate, job)
	})

	t.Run("create jobs and immediately get created jobs", func(t *testing.T) {
		ctx := context.Background()
		db := sqlitedb.NewTestDB()
		companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
		companyService := company.NewDefaultCompanyService(companyRepository)
		jobRepository := sqlitedb.NewSqliteJobRepository(db)
		jobService := job.NewDefaultJobService(jobRepository, companyService)

		jobsToCreate := []jobsummoner.Job{
			{
				Position:    "Software Developer",
				URL:         "https://linkedin.com/jobs/1",
				Location:    "San Francisco",
				CompanyID:   "/google",
				CompanyName: "Google",
				SourceID:    "1",
			},
			{
				Position:    "Manager",
				URL:         "https://linkedin.com/jobs/2",
				Location:    "Seattle",
				CompanyID:   "/microsoft",
				CompanyName: "Microsoft",
				SourceID:    "1",
			},
		}

		ids, errs := jobService.CreateJobs(ctx, jobsToCreate)
		assert.Equal(t, len(jobsToCreate), len(ids))
		assert.Equal(t, 0, len(errs))

		for i, jobToCreate := range jobsToCreate {
			id := ids[i]
			assert.NotEmpty(t, id)

			doesCompanyExist, err := companyService.DoesCompanyExist(ctx, jobToCreate.CompanyID)
			assert.NoError(t, err)
			assert.Equal(t, true, doesCompanyExist)

			job, err := jobService.GetJob(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, jobToCreate, job)
		}
	})
}
