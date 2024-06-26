package main

import (
	"database/sql"
	"log/slog"

	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
)

func newServer(logger *slog.Logger, db *sql.DB) *Server {
	queries := sqlitedb.New(db)

	companyService := company.NewDefaultCompanyService(queries)
	jobService := job.NewDefaultJobService(queries, companyService)
	users := models.UserModel{Queries: queries}

	server := NewServer(logger, jobService, users, db)

	return server
}
