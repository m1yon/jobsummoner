package main

import (
	"net/http"
)

func HomepageServer(w http.ResponseWriter, r *http.Request) {
	jobPostings := GetJobPostings()
	m := NewHomepageViewModel(jobPostings)

	RenderHomepage(m, w)
}
