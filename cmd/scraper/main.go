package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/m1yon/jobsummoner"
)

func main() {
	config := jobsummoner.ScrapeConfig{
		Keywords: []string{"go"},
		Location: "United States",
		MaxAge:   time.Hour * 12,
	}

	url := jobsummoner.BuildURL(config)
	resp, err := http.Get(url)

	if err != nil {
		slog.Error(err.Error())
	}

	results, errs := jobsummoner.ScrapeLinkedInPage(resp.Body)

	if len(errs) != 0 {
		for _, err := range errs {
			slog.Error(err.Error())
		}
	}

	fmt.Println(results)
}
