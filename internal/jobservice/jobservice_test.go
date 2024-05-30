package jobservice

import (
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	jobService := NewDefaultJobService()
	res := jobService.GetJobs()

	assert.Equal(t, []jobsummoner.Job{
		{Position: "Software Engineer"},
		{Position: "Manager"},
	}, res)
}
