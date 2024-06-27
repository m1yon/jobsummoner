package main

import (
	"database/sql"
	"log/slog"

	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/models"
)

func newServer(logger *slog.Logger, db *sql.DB) *Server {
	queries := database.New(db)

	companies := &models.CompanyModel{Queries: queries}
	jobs := &models.JobModel{Queries: queries, Companies: companies}
	users := &models.UserModel{Queries: queries}

	server := NewServer(logger, jobs, users, db)

	return server
}
