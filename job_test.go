package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJobPostings(t *testing.T) {
	jobService := NewDefaultJobService()
	res := jobService.GetJobPostings()

	assert.Equal(t, []JobPosting{
		{"Software Engineer"},
		{"Manager"},
	}, res)
}
