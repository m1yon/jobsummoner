package models

import (
	"time"
)

type ScraperModelInterface interface {
	ScrapeJobs(lastScraped time.Time) ([]Job, []error)
	GetSourceID() string
}
