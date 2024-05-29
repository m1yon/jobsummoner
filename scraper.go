package jobsummoner

import (
	"fmt"

	"github.com/benbjohnson/clock"
	"github.com/go-co-op/gocron/v2"
	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
)

type WorkType string

const (
	WorkTypeOnSite WorkType = "1"
	WorkTypeRemote WorkType = "2"
	WorkTypeHybrid WorkType = "3"
)

type JobType string

const (
	JobTypeFullTime   JobType = "F"
	JobTypePartTime   JobType = "P"
	JobTypeContract   JobType = "C"
	JobTypeTemporary  JobType = "T"
	JobTypeVolunteer  JobType = "V"
	JobTypeInternship JobType = "I"
	JobTypeOther      JobType = "O"
)

type SalaryRange string

const (
	SalaryRange40kPlus  SalaryRange = "1"
	SalaryRange60kPlus  SalaryRange = "2"
	SalaryRange80kPlus  SalaryRange = "3"
	SalaryRange100kPlus SalaryRange = "4"
	SalaryRange120kPlus SalaryRange = "5"
	SalaryRange140kPlus SalaryRange = "6"
	SalaryRange160kPlus SalaryRange = "7"
	SalaryRange180kPlus SalaryRange = "8"
	SalaryRange200kPlus SalaryRange = "9"
)

type ScrapedJob struct {
	Position    string
	CompanyID   string
	CompanyName string
	Location    string
	URL         string
}

type ScrapedJobsResults struct {
	Jobs []ScrapedJob
}

type Scraper interface {
	ScrapeJobs() (ScrapedJobsResults, []error)
}

func ScrapeLoop(c clockwork.Clock, scraper Scraper, interval clock.Duration) error {
	s, err := gocron.NewScheduler(gocron.WithClock(c))

	if err != nil {
		return errors.Wrap(err, "error scheduling cron")
	}

	_, err = s.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(func() {
			scraper.ScrapeJobs()
			fmt.Println("ran")
		}),
	)

	if err != nil {
		return errors.Wrap(err, "error creating new job")
	}

	s.Start()

	return nil
}
