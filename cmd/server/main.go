package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/m1yon/jobsummoner/pkg/http"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	err := godotenv.Load()

	if err != nil {
		logger.Warn("no .env file found", tint.Err(err))
	}

	db, err := openDB(logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	server := initServer(logger, db)
	logger.Info("server started")
	server.ListenAndServe(":3000")
}

func openDB(logger *slog.Logger) (*sql.DB, error) {
	db, err := sqlitedb.NewDB(logger, &sqlitedb.SqlConnectionOpener{})

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initServer(logger *slog.Logger, db *sql.DB) *http.DefaultServer {
	companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
	companyService := company.NewDefaultCompanyService(companyRepository)
	jobRepository := sqlitedb.NewSqliteJobRepository(db)
	jobService := job.NewDefaultJobService(jobRepository, companyService)

	server := http.NewDefaultServer(logger, jobService)

	return server
}
