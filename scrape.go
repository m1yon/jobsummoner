package jobsummoner

import "time"

type Scraper interface {
	ScrapeJobs() ([]Job, []error)
}

type ScrapeService interface {
	Start(scraper []Scraper, crontab string)
}

type Scrape struct {
	ID        int
	SourceID  string
	CreatedAt time.Time
}
