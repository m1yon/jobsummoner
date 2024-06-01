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
	// TODO: add error case when company does not exist
	t.Run("create job and immediately get created job", func(t *testing.T) {
		ctx := context.Background()
		db := sqlitedb.NewTestDB()
		companyRepository := sqlitedb.NewInMemorySqliteCompanyRepository(db)
		companyService := company.NewDefaultCompanyService(companyRepository)
		jobRepository := sqlitedb.NewInMemorySqliteJobRepository(db)
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
}
