package scrape

import (
	"context"
	"log/slog"
	"strings"
	"time"

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

func (ss *DefaultScrapeService) Start(scrapers []jobsummoner.Scraper, crontab string) {
	ctx := context.Background()
	ss.logger.Info("Initializing scrape scheduler...")
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

	_, err = s.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(func() {
			ss.logger.Info("scraping jobs...")

			numberOfJobsScraped := 0

			for _, scraper := range scrapers {
				lastScrapedTime := ss.getLastScrapedTime(ctx, scraper.GetSourceID())
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
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	if err != nil {
		ss.logger.Error(errors.Wrap(err, "error creating new job").Error())
		return
	}

	ss.logger.Info("Scrape scheduler initialized")
	s.Start()

	select {}
}

func (ss *DefaultScrapeService) getLastScrapedTime(ctx context.Context, sourceID string) time.Time {
	defaultLastScrapedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	lastScrape, err := ss.scrapeRepository.GetLastScrape(ctx, sourceID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return lastScrape.CreatedAt
		}

		ss.logger.Error("failed fetching last scraped time")
		return defaultLastScrapedTime
	}

	return defaultLastScrapedTime
}
