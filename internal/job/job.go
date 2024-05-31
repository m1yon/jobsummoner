package job

import (
	"context"

	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type DefaultJobService struct {
	jobRepository jobsummoner.JobRepository
}

func NewDefaultJobService(repository jobsummoner.JobRepository) *DefaultJobService {
	return &DefaultJobService{repository}
}

func (j *DefaultJobService) GetJob(ctx context.Context, id string) (jobsummoner.Job, error) {
	job, err := j.jobRepository.GetJob(ctx, id)

	if err != nil {
		return jobsummoner.Job{}, errors.Wrap(err, "error getting job in job service")
	}

	return job, nil
}

func (j *DefaultJobService) GetJobs(ctx context.Context) ([]jobsummoner.Job, []error) {
	return []jobsummoner.Job{
		{Position: "Software Engineer"},
		{Position: "Manager"},
	}, nil
}

func (j *DefaultJobService) CreateJobs(ctx context.Context, jobs []jobsummoner.Job) []error {
	return nil
}

func (j *DefaultJobService) CreateJob(ctx context.Context, job jobsummoner.Job) (string, error) {
	id, err := j.jobRepository.CreateJob(ctx, job)

	if err != nil {
		return "", errors.Wrap(err, "error creating job in job service")
	}

	return id, nil
}
