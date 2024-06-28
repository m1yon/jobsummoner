package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/pkg/errors"
)

type scrapeService struct {
	logger     *slog.Logger
	httpClient *http.Client
	scrapers   []models.ScraperModelInterface
	c          clockwork.Clock
	scrapes    models.ScrapeModelInterface
	jobs       models.JobModelInterface
}

func newScrapeService(logger *slog.Logger) *scrapeService {
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

	httpClient := newHttpClient(logger, config)

	return &scrapeService{logger: logger, httpClient: httpClient, scrapes: scrapes, jobs: jobs}
}

func (a *scrapeService) addScrapers(scrapers []models.ScraperModelInterface) {
	a.scrapers = scrapers
}

func (ss *scrapeService) start(cron string, scrapeImmediately bool) {
	ctx := context.Background()
	ss.logger.Info("initializing scrape scheduler...")
	s, err := gocron.NewScheduler(gocron.WithClock(ss.c))
	defer (func() {
		err := s.Shutdown()

		if err != nil {
			ss.logger.Error("problem shutting scheduler down", slog.String("err", err.Error()))
			return
		}
	})()

	if err != nil {
		ss.logger.Error(errors.Wrap(err, "error initializing cron scheduler").Error())
		return
	}

	additionalPlugins := make([]gocron.JobOption, 0)

	if scrapeImmediately {
		additionalPlugins = append(additionalPlugins, gocron.WithStartAt(gocron.WithStartImmediately()))
	}

	_, err = s.NewJob(
		gocron.CronJob(cron, false),
		gocron.NewTask(func() {
			ss.logger.Info("scraping jobs...")

			numberOfJobsScraped := 0

			for _, scraper := range ss.scrapers {
				lastScrapedTime, err := ss.scrapes.GetLastScrapeTime(ctx, scraper.GetSourceID())

				if err != nil {
					ss.logger.Error("failed getting last scraped time", slog.String("err", err.Error()))
				}

				results, errs := scraper.ScrapeJobs(lastScrapedTime)

				for _, err := range errs {
					ss.logger.Error("job scrape failure", slog.String("err", err.Error()))
				}

				ss.jobs.CreateJobs(ctx, results)

				err = ss.scrapes.CreateScrape(ctx, scraper.GetSourceID(), ss.c.Now())

				if err != nil {
					ss.logger.Error("failed creating scrape")
				}

				numberOfJobsScraped += len(results)
			}

			ss.logger.Info("scrape successful", slog.Int("jobs", numberOfJobsScraped))
		}),
		additionalPlugins...,
	)

	if err != nil {
		ss.logger.Error(errors.Wrap(err, "error creating new job").Error())
		return
	}

	ss.logger.Info("scrape scheduler initialized")
	s.Start()

	select {}
}
