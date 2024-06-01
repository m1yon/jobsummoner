package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/scrape"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/m1yon/jobsummoner/pkg/linkedin"
	_ "modernc.org/sqlite"
)

func main() {
	url := linkedin.BuildLinkedInJobsURL(linkedin.BuildLinkedInJobsURLArgs{
		Keywords: []string{"go"},
		Location: "United States",
		MaxAge:   time.Hour * 12,
	})
	resp, err := http.Get(url)

	if err != nil {
		slog.Error(err.Error())
	}

	scraper := linkedin.NewLinkedInJobScraper(resp.Body)

	c := clockwork.NewRealClock()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
	scrapeService := scrape.NewDefaultScrapeService(c, logger, jobService)

	scrapeService.Start(scraper, "TZ=America/Denver */30 7-22 * * *")
}
