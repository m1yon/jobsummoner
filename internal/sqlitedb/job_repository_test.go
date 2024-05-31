package sqlitedb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func NewInMemorySqliteJobRepository(db *sql.DB) *SqliteJobRepository {
	queries := New(db)
	return &SqliteJobRepository{queries}
}

func TestJobRepository(t *testing.T) {
	t.Run("add job and immediately get added job", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		jobRepository := NewInMemorySqliteJobRepository(db)

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
