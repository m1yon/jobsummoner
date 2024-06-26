package models

import (
	"time"

	"github.com/m1yon/jobsummoner/internal/sqlitedb"
)

type ScraperModelInterface interface {
	ScrapeJobs(lastScraped time.Time) ([]Job, []error)
	GetSourceID() string
}

type Scraper struct {
}

type ScraperModel struct {
	Queries   *sqlitedb.Queries
	Companies CompanyModelInterface
}
