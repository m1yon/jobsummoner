package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/m1yon/jobsummoner/pkg/http"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := sql.Open("libsql", databaseURL)

	if err != nil {
		logger.Error("failed starting db")
	}

	companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
	companyService := company.NewDefaultCompanyService(companyRepository)
	jobRepository := sqlitedb.NewSqliteJobRepository(db)
	jobService := job.NewDefaultJobService(jobRepository, companyService)

	server := http.NewDefaultServer(logger, jobService)
	server.ListenAndServe(":3000")
}
