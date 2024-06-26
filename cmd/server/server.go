package main

import (
	"database/sql"
	"log/slog"

	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
)

func newServer(logger *slog.Logger, db *sql.DB) *Server {
	queries := sqlitedb.New(db)

	companies := &models.CompanyModel{Queries: queries}
	jobs := &models.JobModel{Queries: queries, Companies: companies}
	users := &models.UserModel{Queries: queries}

	server := NewServer(logger, jobs, users, db)

	return server
}
