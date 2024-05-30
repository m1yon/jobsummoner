package jobsummoner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	jobService := NewDefaultJobService()
	res := jobService.GetJobs()

	assert.Equal(t, []Job{
		{Position: "Software Engineer"},
		{Position: "Manager"},
	}, res)
}
