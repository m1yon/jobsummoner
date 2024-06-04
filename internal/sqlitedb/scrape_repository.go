package sqlitedb

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type SqliteScrapeRepository struct {
	queries *Queries
}

func (s *SqliteScrapeRepository) CreateScrape(ctx context.Context, sourceID string) error {
	err := s.queries.CreateScrape(ctx, sourceID)

	if err != nil {
		return errors.Wrap(err, "problem creating scrape")
	}

	return nil
}

func NewSqliteScrapeRepository(db *sql.DB) *SqliteScrapeRepository {
	queries := New(db)
	return &SqliteScrapeRepository{queries}
}
