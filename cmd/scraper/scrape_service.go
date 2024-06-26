package main

import (
	"database/sql"
	"log/slog"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/m1yon/jobsummoner/internal/scrape"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
)

func newScrapeService(logger *slog.Logger, db *sql.DB) *scrape.DefaultScrapeService {
	c := clockwork.NewRealClock()
	queries := sqlitedb.New(db)

	companies := &models.CompanyModel{Queries: queries}
	jobs := &models.JobModel{Queries: queries, Companies: companies}

	scrapeRepository := models.NewSqliteScrapeRepository(db, c)
	scrapeService := scrape.NewDefaultScrapeService(c, logger, scrapeRepository, jobs)

	return scrapeService
}
