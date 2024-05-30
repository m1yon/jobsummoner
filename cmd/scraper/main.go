package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner"
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
	jobsummoner.ScrapeLoop(c, scraper, "TZ=America/Denver */30 7-22 * * *", logger)
}
