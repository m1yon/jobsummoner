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
	c          clockwork.Clock
	logger     *slog.Logger
	jobService jobsummoner.JobService
}

func NewDefaultScrapeService(c clockwork.Clock, logger *slog.Logger, jobService jobsummoner.JobService) *DefaultScrapeService {
	return &DefaultScrapeService{c, logger, jobService}
}

func (ss *DefaultScrapeService) Start(scraper jobsummoner.Scraper, crontab string) {
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
			results, errs := scraper.ScrapeJobs()

			for _, err := range errs {
				ss.logger.Error("job scrape failure", slog.String("err", err.Error()))
			}

			ss.jobService.CreateJobs(ctx, results.Jobs)

			ss.logger.Info("scrape successful", slog.Int("jobs", len(results.Jobs)))
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
