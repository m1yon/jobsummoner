package main

type JobPosting struct {
	name string
}

type JobPostingStore interface {
	GetJobPostings() []JobPosting
}

func GetJobPostings() []JobPosting {
	return []JobPosting{
		{"Software Engineer"},
		{"Manager"},
	}
}
