package sqlitedb

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type SqliteScrapeRepository struct {
	queries *Queries
	c       clockwork.Clock
}

func NewSqliteScrapeRepository(db *sql.DB, c clockwork.Clock) *SqliteScrapeRepository {
	queries := New(db)
	return &SqliteScrapeRepository{queries, c}
}

func (s *SqliteScrapeRepository) CreateScrape(ctx context.Context, sourceID string, createdAt time.Time) error {
	err := s.queries.CreateScrape(ctx, CreateScrapeParams{SourceID: sourceID, CreatedAt: createdAt.UTC()})

	if err != nil {
		return errors.Wrap(err, "problem creating scrape")
	}

	return nil
}

func (s *SqliteScrapeRepository) GetLastScrape(ctx context.Context, sourceID string) (jobsummoner.Scrape, error) {
	scrape, err := s.queries.GetLastScrape(ctx, sourceID)

	if err != nil {
		return jobsummoner.Scrape{}, errors.Wrap(err, "problem getting scrape")
	}

	return jobsummoner.Scrape{
		ID:        int(scrape.ID),
		SourceID:  scrape.SourceID,
		CreatedAt: scrape.CreatedAt,
	}, nil
}

func (s *SqliteScrapeRepository) GetLastScrapeTime(ctx context.Context, sourceID string) (time.Time, error) {
	defaultLastScrapedTime := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	lastScrape, err := s.queries.GetLastScrape(ctx, sourceID)

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return defaultLastScrapedTime, errors.Wrap(err, "problem getting last scrape")
		}

		return defaultLastScrapedTime, nil
	}

	return lastScrape.CreatedAt, nil
}
