package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/m1yon/jobsummoner/pkg/linkedin"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	scraperApp := newScraperApp(logger)
	scraperApp.AddScrapers([]models.ScraperModelInterface{
		linkedin.NewLinkedInJobScraper(linkedin.NewHttpLinkedInReader(linkedin.LinkedInReaderConfig{
			Keywords: []string{"typescript"},
			Location: "United States",
		}, scraperApp.httpClient, scraperApp.logger), scraperApp.logger),
	})

	scraperApp.Start("TZ=America/Denver */30 7-22 * * *", true)
}
