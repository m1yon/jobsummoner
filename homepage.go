package main

type HomepageViewModel struct {
	text string
}

func NewHomepageViewModel(jobPostings []JobPosting) (m HomepageViewModel) {
	for _, position := range jobPostings {
		m.text += position.name + ","
	}

	return m
}
