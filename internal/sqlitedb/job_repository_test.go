package sqlitedb

import (
	"context"
	"testing"
	"time"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestJobRepository(t *testing.T) {
	t.Run("create job and immediately get created job", func(t *testing.T) {
		ctx := context.Background()
		db, _ := NewInMemoryDB()
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

		want := jobsummoner.Job{
			Position:      "Software Developer",
			URL:           "https://linkedin.com/jobs/1",
			Location:      "San Francisco",
			LastPosted:    time.Now(),
			SourceID:      companyToCreate.SourceID,
			CompanyID:     companyToCreate.ID,
			CompanyName:   companyToCreate.Name,
			CompanyAvatar: companyToCreate.Avatar,
			CompanyURL:    companyToCreate.Url,
		}

		createdJobID, err := jobRepository.CreateJob(ctx, want)
		assert.NoError(t, err)

		got, err := jobRepository.GetJob(ctx, createdJobID)
		assert.NoError(t, err)
		assertJobsEqual(t, want, got)
	})

	t.Run("get jobs", func(t *testing.T) {
		ctx := context.Background()
		db, _ := NewInMemoryDB()
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
				LastPosted:    time.Now(),
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
				LastPosted:    time.Now(),
			},
		}

		_, _ = jobRepository.CreateJob(ctx, jobsToCreate[0])
		_, _ = jobRepository.CreateJob(ctx, jobsToCreate[1])

		jobs, err := jobRepository.GetJobs(ctx)
		assert.NoError(t, err)
		assertJobListsEqual(t, jobsToCreate, jobs)
	})
}

func assertJobListsEqual(t *testing.T, expectedJobList, actualJobList []jobsummoner.Job) {
	t.Helper()

	if len(expectedJobList) != len(actualJobList) {
		t.Fatalf("expected %v jobs, got %v", len(expectedJobList), len(actualJobList))
	}

	for i := range expectedJobList {
		assertJobsEqual(t, expectedJobList[i], actualJobList[i])
	}
}

func assertJobsEqual(t *testing.T, expectedJob, actualJob jobsummoner.Job) {
	t.Helper()

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
