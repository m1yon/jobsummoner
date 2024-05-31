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

func NewSqliteJobRepository(dataSourceName string) *SqliteJobRepository {
	db, _ := sql.Open("sqlite", dataSourceName)
	queries := New(db)
	return &SqliteJobRepository{queries}
}
