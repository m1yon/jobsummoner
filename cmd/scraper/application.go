package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/models"
)

type scraperApp struct {
	logger        *slog.Logger
	db            *sql.DB
	scrapeService *ScrapeService
	httpClient    *http.Client
	scrapers      []models.ScraperModelInterface
	config        *config
}

func newScraperApp(logger *slog.Logger) *scraperApp {
	config := getConfigFromFlags()

	db, err := openDB(logger, config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	c := clockwork.NewRealClock()
	queries := database.New(db)

	companies := &models.CompanyModel{Queries: queries}
	jobs := &models.JobModel{Queries: queries, Companies: companies}
	scrapes := &models.ScrapeModel{Queries: queries, C: c}

	scrapeService := &ScrapeService{c: c, logger: logger, scrapes: scrapes, jobs: jobs}

	httpClient := newHttpClient(logger, config)

	return &scraperApp{logger: logger, db: db, scrapeService: scrapeService, httpClient: httpClient, config: config}
}

func (a *scraperApp) AddScrapers(scrapers []models.ScraperModelInterface) {
	a.scrapers = scrapers
}

func (a *scraperApp) Start(cron string, scrapeImmediately bool) {
	a.logger.Info("scraper app started")
	a.scrapeService.Start(a.scrapers, cron, scrapeImmediately)
}
