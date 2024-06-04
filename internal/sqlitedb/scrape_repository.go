package sqlitedb

import (
	"context"
	"database/sql"
	"time"

	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type SqliteScrapeRepository struct {
	queries *Queries
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

func NewSqliteScrapeRepository(db *sql.DB) *SqliteScrapeRepository {
	queries := New(db)
	return &SqliteScrapeRepository{queries}
}
