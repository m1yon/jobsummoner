package models

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/pkg/errors"
)

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
	GetLastScrapeTime(ctx context.Context, sourceID string) (time.Time, error)
}

type SqliteScrapeRepository struct {
	queries *sqlitedb.Queries
	c       clockwork.Clock
}

func NewSqliteScrapeRepository(db *sql.DB, c clockwork.Clock) *SqliteScrapeRepository {
	queries := sqlitedb.New(db)
	return &SqliteScrapeRepository{queries, c}
}

func (s *SqliteScrapeRepository) CreateScrape(ctx context.Context, sourceID string, createdAt time.Time) error {
	err := s.queries.CreateScrape(ctx, sqlitedb.CreateScrapeParams{SourceID: sourceID, CreatedAt: createdAt.UTC()})

	if err != nil {
		return errors.Wrap(err, "problem creating scrape")
	}

	return nil
}

func (s *SqliteScrapeRepository) GetLastScrape(ctx context.Context, sourceID string) (Scrape, error) {
	scrape, err := s.queries.GetLastScrape(ctx, sourceID)

	if err != nil {
		return Scrape{}, errors.Wrap(err, "problem getting scrape")
	}

	return Scrape{
		ID:        int(scrape.ID),
		SourceID:  scrape.SourceID,
		CreatedAt: scrape.CreatedAt,
	}, nil
}

func (s *SqliteScrapeRepository) GetLastScrapeTime(ctx context.Context, sourceID string) (time.Time, error) {
	defaultLastScrapedTime := s.c.Now().Add(-24 * time.Hour)
	lastScrape, err := s.queries.GetLastScrape(ctx, sourceID)

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return defaultLastScrapedTime, errors.Wrap(err, "problem getting last scrape")
		}

		return defaultLastScrapedTime, nil
	}

	return lastScrape.CreatedAt, nil
}
