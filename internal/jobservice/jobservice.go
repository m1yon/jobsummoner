package jobservice

import "github.com/m1yon/jobsummoner"

type DefaultJobService struct{}

func NewDefaultJobService() *DefaultJobService {
	return &DefaultJobService{}
}

func (j *DefaultJobService) GetJobs() []jobsummoner.Job {
	return []jobsummoner.Job{
		{Position: "Software Engineer"},
		{Position: "Manager"},
	}
}
