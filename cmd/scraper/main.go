package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/scrape"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/m1yon/jobsummoner/pkg/linkedin"
	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	scrapers := []jobsummoner.Scraper{
		linkedin.NewLinkedInJobScraper(linkedin.LinkedInReaderConfig{
			Keywords: []string{"typescript"},
			Location: "United States",
		}, logger),
	}

	c := clockwork.NewRealClock()

	db, err := sql.Open("sqlite", "./db/database.db")

	if err != nil {
		logger.Error("failed starting db")
	}

	err = db.Ping()

	if err != nil {
		logger.Error("failed pinging db")
	}

	companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
	companyService := company.NewDefaultCompanyService(companyRepository)
	jobRepository := sqlitedb.NewSqliteJobRepository(db)
	jobService := job.NewDefaultJobService(jobRepository, companyService)

	scrapeRepository := sqlitedb.NewSqliteScrapeRepository(db)
	scrapeService := scrape.NewDefaultScrapeService(c, logger, scrapeRepository, jobService)

	scrapeService.Start(scrapers, "TZ=America/Denver */30 7-22 * * *")
}
