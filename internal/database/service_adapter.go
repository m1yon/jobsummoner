package database

import (
	"context"
	"database/sql"

	"github.com/m1yon/jobsummoner"
)

type SQLCAdapter struct {
	queries *Queries
}

func (s *SQLCAdapter) AddJob(ctx context.Context, arg jobsummoner.AddJobParams) error {
	return nil
}

func NewSQLCAdapter() *SQLCAdapter {
	db, _ := sql.Open("sqlite", "./db/database.db")
	queries := New(db)
	return &SQLCAdapter{queries}
}
