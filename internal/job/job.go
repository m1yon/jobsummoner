package job

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/pkg/errors"
)

type DefaultJobService struct {
	Queries   *sqlitedb.Queries
	companies models.CompanyModelInterface
}

func NewDefaultJobService(queries *sqlitedb.Queries, companies models.CompanyModelInterface) *DefaultJobService {
	return &DefaultJobService{queries, companies}
}

func (j *DefaultJobService) GetJob(ctx context.Context, id string) (jobsummoner.Job, error) {
	dbJob, err := j.Queries.GetJob(ctx, id)

	if err != nil {
		return jobsummoner.Job{}, errors.Wrap(err, "error getting job from db")
	}

	job := jobsummoner.Job{
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

func (j *DefaultJobService) GetJobs(ctx context.Context) ([]jobsummoner.Job, error) {
	jobs, err := j.Queries.GetJobs(ctx)

	if err != nil {
		return []jobsummoner.Job{}, errors.Wrap(err, "error getting jobs from db")
	}

	formattedJobs := make([]jobsummoner.Job, 0, len(jobs))

	for _, dbJob := range jobs {
		formattedJobs = append(formattedJobs, jobsummoner.Job{
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
	doesCompanyExist, err := j.companies.DoesCompanyExist(ctx, job.CompanyID)

	if err != nil {
		return "", errors.Wrap(err, "error fetching company in job service")
	}

	if !doesCompanyExist {
		_, err := j.companies.CreateCompany(ctx, models.Company{ID: job.CompanyID, Name: job.CompanyName, SourceID: job.SourceID, Url: job.CompanyURL, Avatar: job.CompanyAvatar})

		if err != nil {
			return "", errors.Wrap(err, "error creating company in job service")
		}
	}

	id := generateJobID(job.CompanyID, job.Position)
	err = j.Queries.CreateJob(ctx, sqlitedb.CreateJobParams{
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
