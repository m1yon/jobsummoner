package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/models"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	config := getConfigFromFlags()

	db, err := openDB(logger, config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	queries := database.New(db)

	companies := &models.CompanyModel{Queries: queries}
	jobs := &models.JobModel{Queries: queries, Companies: companies}
	users := &models.UserModel{Queries: queries}

	server := newServer(logger, jobs, users, db)

	server.Start(":3000")
}
