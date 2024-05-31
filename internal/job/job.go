package job

import (
	"github.com/m1yon/jobsummoner"
)

type DefaultJobService struct {
	jobRepository jobsummoner.JobRepository
}

func NewDefaultJobService(repository jobsummoner.JobRepository) *DefaultJobService {
	return &DefaultJobService{repository}
}

func (j *DefaultJobService) GetJobs() []jobsummoner.Job {
	return []jobsummoner.Job{
		{Position: "Software Engineer"},
		{Position: "Manager"},
	}
}

func (j *DefaultJobService) AddJobs(jobs []jobsummoner.Job) {

}
