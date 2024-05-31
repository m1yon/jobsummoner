package job

import (
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	jobRepository := sqlitedb.NewSqliteJobRepository()
	jobService := NewDefaultJobService(jobRepository)

	res := jobService.GetJobs()

	assert.Equal(t, []jobsummoner.Job{
		{Position: "Software Engineer"},
		{Position: "Manager"},
	}, res)
}

func TestAddJobs(t *testing.T) {
	jobRepository := sqlitedb.NewSqliteJobRepository()
	jobService := NewDefaultJobService(jobRepository)

	jobs := []jobsummoner.Job{
		{Position: "Software Engineer"},
		{Position: "Manager"},
	}

	jobService.AddJobs(jobs)

	// assert.Equal(t, []jobsummoner.Job{
	// 	{Position: "Software Engineer"},
	// 	{Position: "Manager"},
	// }, res)
}
