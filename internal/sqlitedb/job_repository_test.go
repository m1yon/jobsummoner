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
		jobRepository := NewInMemorySqliteJobRepository()

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
}
