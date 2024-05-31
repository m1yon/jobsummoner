package jobservice

import (
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	jobServiceStore := database.NewSQLCAdapter()
	jobService := NewDefaultJobService(jobServiceStore)

	res := jobService.GetJobs()

	assert.Equal(t, []jobsummoner.Job{
		{Position: "Software Engineer"},
		{Position: "Manager"},
	}, res)
}

func TestAddJobs(t *testing.T) {
	jobServiceStore := database.NewSQLCAdapter()
	jobService := NewDefaultJobService(jobServiceStore)

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
