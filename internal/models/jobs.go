package models

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/pkg/errors"
)

type JobModelInterface interface {
	GetJob(ctx context.Context, id string) (Job, error)
	GetJobs(ctx context.Context) ([]Job, error)
	CreateJob(ctx context.Context, job Job) (string, error)
	CreateJobs(ctx context.Context, jobs []Job) ([]string, []error)
}

type Job struct {
	Position      string
	Location      string
	URL           string
	SourceID      string
	LastPosted    time.Time
	CompanyID     string
	CompanyName   string
	CompanyAvatar string
	CompanyURL    string
}

type JobModel struct {
	Queries   *sqlitedb.Queries
	Companies CompanyModelInterface
}

func (m *JobModel) GetJob(ctx context.Context, id string) (Job, error) {
	dbJob, err := m.Queries.GetJob(ctx, id)

	if err != nil {
		return Job{}, errors.Wrap(err, "error getting job from db")
	}

	job := Job{
		Position:      dbJob.Position,
		Location:      dbJob.Location.String,
		URL:           dbJob.JobUrl,
		SourceID:      dbJob.SourceID,
		LastPosted:    dbJob.LastPosted,
		CompanyID:     dbJob.CompanyID,
		CompanyName:   dbJob.CompanyName,
		CompanyAvatar: dbJob.CompanyAvatar.String,
		CompanyURL:    dbJob.CompanyUrl,
	}

	return job, nil
}

func (m *JobModel) GetJobs(ctx context.Context) ([]Job, error) {
	jobs, err := m.Queries.GetJobs(ctx)

	if err != nil {
		return []Job{}, errors.Wrap(err, "error getting jobs from db")
	}

	formattedJobs := make([]Job, 0, len(jobs))

	for _, dbJob := range jobs {
		formattedJobs = append(formattedJobs, Job{
			Position:      dbJob.Position,
			Location:      dbJob.Location.String,
			URL:           dbJob.JobUrl,
			SourceID:      dbJob.SourceID,
			LastPosted:    dbJob.LastPosted,
			CompanyID:     dbJob.CompanyID,
			CompanyName:   dbJob.CompanyName,
			CompanyAvatar: dbJob.CompanyAvatar.String,
			CompanyURL:    dbJob.CompanyUrl,
		})
	}

	return formattedJobs, nil
}

func (m *JobModel) CreateJobs(ctx context.Context, jobs []Job) ([]string, []error) {
	ids := make([]string, 0, len(jobs))
	errs := make([]error, 0)

	for _, job := range jobs {
		id, err := m.CreateJob(ctx, job)

		if err != nil {
			errs = append(errs, err)
		}

		ids = append(ids, id)
	}

	return ids, errs
}

func (m *JobModel) CreateJob(ctx context.Context, job Job) (string, error) {
	doesCompanyExist, err := m.Companies.DoesCompanyExist(ctx, job.CompanyID)

	if err != nil {
		return "", errors.Wrap(err, "error fetching company in job service")
	}

	if !doesCompanyExist {
		_, err := m.Companies.CreateCompany(ctx, Company{ID: job.CompanyID, Name: job.CompanyName, SourceID: job.SourceID, Url: job.CompanyURL, Avatar: job.CompanyAvatar})

		if err != nil {
			return "", errors.Wrap(err, "error creating company in job service")
		}
	}

	id := generateJobID(job.CompanyID, job.Position)
	err = m.Queries.CreateJob(ctx, sqlitedb.CreateJobParams{
		ID:       id,
		Position: job.Position,
		Location: sql.NullString{
			Valid:  job.Location != "",
			String: job.Location,
		},
		Url:        job.URL,
		CompanyID:  job.CompanyID,
		SourceID:   job.SourceID,
		LastPosted: job.LastPosted,
	})

	if err != nil {
		return "", errors.Wrap(err, "error creating job in job service")
	}

	return id, nil
}

func generateJobID(company_id string, position string) string {
	data := company_id + "|" + position

	hasher := sha256.New()

	hasher.Write([]byte(data))

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}

type WorkType string

const (
	WorkTypeOnSite WorkType = "1"
	WorkTypeRemote WorkType = "2"
	WorkTypeHybrid WorkType = "3"
)

type JobType string

const (
	JobTypeFullTime   JobType = "F"
	JobTypePartTime   JobType = "P"
	JobTypeContract   JobType = "C"
	JobTypeTemporary  JobType = "T"
	JobTypeVolunteer  JobType = "V"
	JobTypeInternship JobType = "I"
	JobTypeOther      JobType = "O"
)

type SalaryRange string

const (
	SalaryRange40kPlus  SalaryRange = "1"
	SalaryRange60kPlus  SalaryRange = "2"
	SalaryRange80kPlus  SalaryRange = "3"
	SalaryRange100kPlus SalaryRange = "4"
	SalaryRange120kPlus SalaryRange = "5"
	SalaryRange140kPlus SalaryRange = "6"
	SalaryRange160kPlus SalaryRange = "7"
	SalaryRange180kPlus SalaryRange = "8"
	SalaryRange200kPlus SalaryRange = "9"
)
