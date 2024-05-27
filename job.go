package main

type JobService interface {
	GetJobPostings() []JobPosting
}

type DefaultJobService struct{}

func NewDefaultJobService() *DefaultJobService {
	return &DefaultJobService{}
}

type JobPosting struct {
	name string
}

type JobPostingStore interface {
	GetJobPostings() []JobPosting
}

func (j *DefaultJobService) GetJobPostings() []JobPosting {
	return []JobPosting{
		{"Software Engineer"},
		{"Manager"},
	}
}
