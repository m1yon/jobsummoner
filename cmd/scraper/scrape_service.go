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

	companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
	companyService := company.NewDefaultCompanyService(companyRepository)
	jobRepository := sqlitedb.NewSqliteJobRepository(db)
	jobService := job.NewDefaultJobService(jobRepository, companyService)

	scrapeRepository := sqlitedb.NewSqliteScrapeRepository(db, c)
	scrapeService := scrape.NewDefaultScrapeService(c, logger, scrapeRepository, jobService)

	return scrapeService
}
