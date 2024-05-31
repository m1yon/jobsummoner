package sqlitedb

import (
	"context"
	"database/sql"

	"github.com/m1yon/jobsummoner"
)

type SqliteJobRepository struct {
	queries *Queries
}

func (s *SqliteJobRepository) AddJob(ctx context.Context, arg jobsummoner.AddJobParams) error {
	return nil
}

func NewSqliteJobRepository() *SqliteJobRepository {
	db, _ := sql.Open("sqlite", "./db/database.db")
	queries := New(db)
	return &SqliteJobRepository{queries}
}
