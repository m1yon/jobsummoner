package sqlitedb

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"

	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type SqliteJobRepository struct {
	queries *Queries
}

func (s *SqliteJobRepository) CreateJob(ctx context.Context, arg jobsummoner.Job) (string, error) {
	id := generateJobID(arg.CompanyID, arg.Position)
	err := s.queries.CreateJob(ctx, CreateJobParams{
		ID:       id,
		Position: arg.Position,
		Location: sql.NullString{
			Valid:  arg.Location != "",
			String: arg.Location,
		},
		Url:        arg.URL,
		CompanyID:  arg.CompanyID,
		SourceID:   arg.SourceID,
		LastPosted: arg.LastPosted,
	})

	if err != nil {
		return "", errors.Wrap(err, "error creating job in db")
	}

	return id, nil
}

func (s *SqliteJobRepository) GetJob(ctx context.Context, id string) (jobsummoner.Job, error) {
	dbJob, err := s.queries.GetJob(ctx, id)

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

func (s *SqliteJobRepository) GetJobs(ctx context.Context) ([]jobsummoner.Job, error) {
	jobs, err := s.queries.GetJobs(ctx)

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

func NewSqliteJobRepository(db *sql.DB) *SqliteJobRepository {
	queries := New(db)
	return &SqliteJobRepository{queries}
}

func generateJobID(company_id string, position string) string {
	data := company_id + "|" + position

	hasher := sha256.New()

	hasher.Write([]byte(data))

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}
