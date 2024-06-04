package jobsummoner

import (
	"context"
	"time"
)

type Scraper interface {
	ScrapeJobs(lastScraped time.Time) ([]Job, []error)
	GetSourceID() string
}

type ScrapeService interface {
	Start(scraper []Scraper, crontab string)
}

type Scrape struct {
	ID        int
	SourceID  string
	CreatedAt time.Time
}

type ScrapeRepository interface {
	CreateScrape(ctx context.Context, sourceID string, createdAt time.Time) error
	GetLastScrape(ctx context.Context, sourceID string) (Scrape, error)
}
