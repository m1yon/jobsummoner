package main

import (
	"context"
	"net/http"
)

func HomepageServer(w http.ResponseWriter, r *http.Request) {
	jobPostings := GetJobPostings()
	m := NewHomepageViewModel(jobPostings)

	component := homepage(m)
	component.Render(context.Background(), w)
}

type HomepageViewModel struct {
	text string
}

func NewHomepageViewModel(jobPostings []JobPosting) (m HomepageViewModel) {
	for _, position := range jobPostings {
		m.text += position.name + ","
	}

	return m
}
