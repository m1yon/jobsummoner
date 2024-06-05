package scrape

import (
	"context"
	"log/slog"

	"github.com/go-co-op/gocron/v2"
	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type DefaultScrapeService struct {
	c                clockwork.Clock
	logger           *slog.Logger
	scrapeRepository jobsummoner.ScrapeRepository
	jobService       jobsummoner.JobService
}

func NewDefaultScrapeService(c clockwork.Clock, logger *slog.Logger, repository jobsummoner.ScrapeRepository, jobService jobsummoner.JobService) *DefaultScrapeService {
	return &DefaultScrapeService{c, logger, repository, jobService}
}

func (ss *DefaultScrapeService) Start(scrapers []jobsummoner.Scraper, crontab string, scrapeImmediately bool) {
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
		gocron.CronJob(crontab, false),
		gocron.NewTask(func() {
			ss.logger.Info("scraping jobs...")

			numberOfJobsScraped := 0

			for _, scraper := range scrapers {
				lastScrapedTime, err := ss.scrapeRepository.GetLastScrapeTime(ctx, scraper.GetSourceID())

				if err != nil {
					ss.logger.Error("failed getting last scraped time", slog.String("err", err.Error()))
				}

				results, errs := scraper.ScrapeJobs(lastScrapedTime)

				for _, err := range errs {
					ss.logger.Error("job scrape failure", slog.String("err", err.Error()))
				}

				ss.jobService.CreateJobs(ctx, results)

				err = ss.scrapeRepository.CreateScrape(ctx, scraper.GetSourceID(), ss.c.Now())

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
