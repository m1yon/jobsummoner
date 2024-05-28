package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

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
	results, errs := scraper.ScrapeJobs()

	if len(errs) != 0 {
		for _, err := range errs {
			slog.Error(err.Error())
		}
	}

	fmt.Println(results)
}
