package main

import (
	"net/http"
)

func HomepageServer(w http.ResponseWriter, r *http.Request) {
	jobPostings := GetJobPostings()
	m := NewHomepageViewModel(jobPostings)

	w.Write([]byte(m.text))
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
