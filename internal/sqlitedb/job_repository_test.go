package sqlitedb

import (
	"context"
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestJobRepository(t *testing.T) {
	t.Run("create job and immediately get created job", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		companyRepository := NewSqliteCompanyRepository(db)
		jobRepository := NewSqliteJobRepository(db)

		companyToCreate := jobsummoner.Company{
			ID:       "/google",
			Name:     "Google",
			Url:      "https://google.com/",
			Avatar:   "https://upload.wikimedia.org/wikipedia/commons/c/c1/Google_%22G%22_logo.svg",
			SourceID: "linkedin",
		}

		id, err := companyRepository.CreateCompany(ctx, companyToCreate)
		assert.NoError(t, err)
		assert.Equal(t, companyToCreate.ID, id)

		jobToCreate := jobsummoner.Job{
			Position:      "Software Developer",
			URL:           "https://linkedin.com/jobs/1",
			Location:      "San Francisco",
			SourceID:      companyToCreate.SourceID,
			CompanyID:     companyToCreate.ID,
			CompanyName:   companyToCreate.Name,
			CompanyAvatar: companyToCreate.Avatar,
			CompanyURL:    companyToCreate.Url,
		}

		createdJobID, err := jobRepository.CreateJob(ctx, jobToCreate)
		assert.NoError(t, err)

		job, err := jobRepository.GetJob(ctx, createdJobID)
		assert.NoError(t, err)
		assert.Equal(t, jobToCreate, job)
	})

	t.Run("get jobs", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		companyRepository := NewSqliteCompanyRepository(db)
		jobRepository := NewSqliteJobRepository(db)

		companyToCreate := jobsummoner.Company{
			ID:       "/google",
			Name:     "Google",
			Url:      "https://google.com/",
			Avatar:   "https://upload.wikimedia.org/wikipedia/commons/c/c1/Google_%22G%22_logo.svg",
			SourceID: "linkedin",
		}

		id, err := companyRepository.CreateCompany(ctx, companyToCreate)
		assert.NoError(t, err)
		assert.Equal(t, companyToCreate.ID, id)

		jobsToCreate := []jobsummoner.Job{
			{
				Position:      "Software Developer",
				URL:           "https://linkedin.com/jobs/1",
				Location:      "San Francisco",
				SourceID:      companyToCreate.SourceID,
				CompanyID:     companyToCreate.ID,
				CompanyName:   companyToCreate.Name,
				CompanyAvatar: companyToCreate.Avatar,
				CompanyURL:    companyToCreate.Url,
			},
			{
				Position:      "Manager",
				URL:           "https://linkedin.com/jobs/2",
				Location:      "San Francisco",
				SourceID:      companyToCreate.SourceID,
				CompanyID:     companyToCreate.ID,
				CompanyName:   companyToCreate.Name,
				CompanyAvatar: companyToCreate.Avatar,
				CompanyURL:    companyToCreate.Url,
			},
		}

		_, _ = jobRepository.CreateJob(ctx, jobsToCreate[0])
		_, _ = jobRepository.CreateJob(ctx, jobsToCreate[1])

		jobs, err := jobRepository.GetJobs(ctx)
		assert.NoError(t, err)
		assert.Equal(t, jobsToCreate, jobs)
	})
}
