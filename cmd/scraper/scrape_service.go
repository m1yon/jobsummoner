package main

import (
	"database/sql"
	"log/slog"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/scrape"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
)

func newScrapeService(logger *slog.Logger, db *sql.DB) *scrape.DefaultScrapeService {
	c := clockwork.NewRealClock()
	queries := sqlitedb.New(db)

	companyService := company.NewDefaultCompanyService(queries)
	jobService := job.NewDefaultJobService(queries, companyService)

	scrapeRepository := sqlitedb.NewSqliteScrapeRepository(db, c)
	scrapeService := scrape.NewDefaultScrapeService(c, logger, scrapeRepository, jobService)

	return scrapeService
}
