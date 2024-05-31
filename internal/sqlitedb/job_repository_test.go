package sqlitedb

import (
	"context"
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestJobRepository(t *testing.T) {
	t.Run("add job and immediately get added job", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		companyRepository := NewInMemorySqliteCompanyRepository(db)
		jobRepository := NewInMemorySqliteJobRepository(db, companyRepository)

		jobToAdd := jobsummoner.Job{
			Position: "Software Developer",
			URL:      "https://linkedin.com/jobs/1",
			Location: "San Francisco",
		}
		id, err := jobRepository.AddJob(ctx, jobToAdd)

		assert.NoError(t, err)

		job, err := jobRepository.GetJob(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, jobToAdd, job)
	})

	t.Run("add job from new company and ensure company is created", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		companyRepository := NewInMemorySqliteCompanyRepository(db)
		jobRepository := NewInMemorySqliteJobRepository(db, companyRepository)

		jobToAdd := jobsummoner.Job{
			Position:  "Software Developer",
			URL:       "https://linkedin.com/jobs/1",
			Location:  "San Francisco",
			CompanyID: "/google",
		}

		doesCompanyExist, err := companyRepository.DoesCompanyExist(ctx, jobToAdd.CompanyName)
		assert.NoError(t, err)
		assert.Equal(t, false, doesCompanyExist, "company shouldn't exist yet")

		id, err := jobRepository.AddJob(ctx, jobToAdd)
		assert.NoError(t, err)

		job, err := jobRepository.GetJob(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, jobToAdd, job)

		doesCompanyExist, err = companyRepository.DoesCompanyExist(ctx, jobToAdd.CompanyID)
		assert.NoError(t, err)
		assert.Equal(t, true, doesCompanyExist, "company should exist now")
	})
}
