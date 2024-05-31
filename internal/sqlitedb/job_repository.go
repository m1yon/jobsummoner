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

func (s *SqliteJobRepository) AddJob(ctx context.Context, arg jobsummoner.Job) (string, error) {
	id := generateJobID(arg.CompanyID, arg.Position)

	err := s.queries.AddJob(ctx, AddJobParams{
		ID:       id,
		Position: arg.Position,
		Location: sql.NullString{
			Valid:  arg.Location != "",
			String: arg.Location,
		},
		Url: arg.URL,
	})

	if err != nil {
		return "", errors.Wrap(err, "error adding job to db")
	}

	return id, nil
}

func (s *SqliteJobRepository) GetJob(ctx context.Context, id string) (jobsummoner.Job, error) {
	dbJob, err := s.queries.GetJob(ctx, id)

	if err != nil {
		return jobsummoner.Job{}, errors.Wrap(err, "error getting job from db")
	}

	job := jobsummoner.Job{
		Position:    dbJob.Position,
		CompanyID:   dbJob.CompanyID,
		CompanyName: "",
		Location:    dbJob.Location.String,
		URL:         dbJob.Url,
	}

	return job, nil
}

func NewSqliteJobRepository(dataSourceName string) *SqliteJobRepository {
	db, _ := sql.Open("sqlite", dataSourceName)
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