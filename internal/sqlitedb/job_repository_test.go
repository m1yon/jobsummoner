package sqlitedb

import (
	"context"
	"fmt"
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestJobRepository(t *testing.T) {
	t.Run("create job and immediately get created job", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		jobRepository := NewInMemorySqliteJobRepository(db)

		jobToCreate := jobsummoner.Job{
			Position: "Software Developer",
			URL:      "https://linkedin.com/jobs/1",
			Location: "San Francisco",
		}
		id, err := jobRepository.CreateJob(ctx, jobToCreate)
		fmt.Println(id)

		assert.NoError(t, err)

		// TODO: assert the rest
		// _, err = jobRepository.GetJob(ctx, id)

		// assert.NoError(t, err)
		// assert.Equal(t, jobToCreate, job)
	})
}
