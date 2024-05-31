package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/scrape"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/m1yon/jobsummoner/pkg/linkedin"
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
	jobRepository := sqlitedb.NewSqliteJobRepository("./db/database.db")
	jobService := job.NewDefaultJobService(jobRepository)
	scrapeService := scrape.NewDefaultScrapeService(c, logger, jobService)
	scrapeService.Start(scraper, "TZ=America/Denver */30 7-22 * * *")
}
