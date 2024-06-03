package job

import (
	"context"

	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type DefaultJobService struct {
	jobRepository  jobsummoner.JobRepository
	companyService jobsummoner.CompanyService
}

func NewDefaultJobService(repository jobsummoner.JobRepository, companyService jobsummoner.CompanyService) *DefaultJobService {
	return &DefaultJobService{repository, companyService}
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

func (j *DefaultJobService) CreateJobs(ctx context.Context, jobs []jobsummoner.Job) ([]string, []error) {
	ids := make([]string, 0, len(jobs))
	errs := make([]error, 0)

	for _, job := range jobs {
		id, err := j.CreateJob(ctx, job)

		if err != nil {
			errs = append(errs, err)
		}

		ids = append(ids, id)
	}

	return ids, errs
}

func (j *DefaultJobService) CreateJob(ctx context.Context, job jobsummoner.Job) (string, error) {
	doesCompanyExist, err := j.companyService.DoesCompanyExist(ctx, job.CompanyID)

	if err != nil {
		return "", errors.Wrap(err, "error fetching company in job service")
	}

	if !doesCompanyExist {
		_, err := j.companyService.CreateCompany(ctx, jobsummoner.Company{ID: job.CompanyID, Name: job.CompanyName})

		if err != nil {
			return "", errors.Wrap(err, "error creating company in job service")
		}
	}

	id, err := j.jobRepository.CreateJob(ctx, job)

	if err != nil {
		return "", errors.Wrap(err, "error creating job in job service")
	}

	return id, nil
}
