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
	companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
	companyService := company.NewDefaultCompanyService(companyRepository)
	jobRepository := sqlitedb.NewSqliteJobRepository(db)
	jobService := job.NewDefaultJobService(jobRepository, companyService)
	userService := models.NewDefaultUserService(db)

	server := NewServer(logger, jobService, userService, db)

	return server
}
