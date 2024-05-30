package jobsummoner

import (
	"log/slog"

	"github.com/go-co-op/gocron/v2"
	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
)

type ScrapedJobsResults struct {
	Jobs []Job
}

type Scraper interface {
	ScrapeJobs() (ScrapedJobsResults, []error)
}

type ScrapeService interface {
	Start(scraper Scraper, crontab string)
}

type DefaultScrapeService struct {
	c      clockwork.Clock
	logger *slog.Logger
}

func NewDefaultScrapeService(c clockwork.Clock, logger *slog.Logger) *DefaultScrapeService {
	return &DefaultScrapeService{c, logger}
}

func (ss *DefaultScrapeService) Start(scraper Scraper, crontab string) {
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

			ss.logger.Info("scrape successful", slog.Int("jobs", len(results.Jobs)))
		}),
	)

	if err != nil {
		ss.logger.Error(errors.Wrap(err, "error creating new job").Error())
		return
	}

	ss.logger.Info("Scrape scheduler initialized")
	s.Start()

	select {}
}
