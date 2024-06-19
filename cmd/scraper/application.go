package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/scrape"
)

type scraperApp struct {
	logger        *slog.Logger
	db            *sql.DB
	scrapeService *scrape.DefaultScrapeService
	httpClient    *http.Client
	scrapers      []jobsummoner.Scraper
	config        *config
}

func newScraperApp(logger *slog.Logger) *scraperApp {
	config := getConfigFromFlags()

	db, err := openDB(logger, config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	scrapeService := newScrapeService(logger, db)
	httpClient := newHttpClient(logger, config)

	return &scraperApp{logger: logger, db: db, scrapeService: scrapeService, httpClient: httpClient, config: config}
}

func (a *scraperApp) AddScrapers(scrapers []jobsummoner.Scraper) {
	a.scrapers = scrapers
}

func (a *scraperApp) Start(cron string, scrapeImmediately bool) {
	a.logger.Info("scraper app started")
	a.scrapeService.Start(a.scrapers, cron, scrapeImmediately)
}
