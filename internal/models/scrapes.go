package models

import (
	"context"
	"strings"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/pkg/errors"
)

type ScrapeModelInterface interface {
	Create(ctx context.Context, sourceID string, createdAt time.Time) error
	Latest(ctx context.Context, sourceID string) (Scrape, error)
	LastRan(ctx context.Context, sourceID string) (time.Time, error)
}

type Scrape struct {
	ID        int
	SourceID  string
	CreatedAt time.Time
}

type ScrapeModel struct {
	Queries *database.Queries
	C       clockwork.Clock
}

func (m *ScrapeModel) Create(ctx context.Context, sourceID string, createdAt time.Time) error {
	err := m.Queries.CreateScrape(ctx, database.CreateScrapeParams{SourceID: sourceID, CreatedAt: createdAt.UTC()})

	if err != nil {
		return errors.Wrap(err, "problem creating scrape")
	}

	return nil
}

func (m *ScrapeModel) Latest(ctx context.Context, sourceID string) (Scrape, error) {
	scrape, err := m.Queries.GetLastScrape(ctx, sourceID)

	if err != nil {
		return Scrape{}, errors.Wrap(err, "problem getting scrape")
	}

	return Scrape{
		ID:        int(scrape.ID),
		SourceID:  scrape.SourceID,
		CreatedAt: scrape.CreatedAt,
	}, nil
}

func (m *ScrapeModel) LastRan(ctx context.Context, sourceID string) (time.Time, error) {
	defaultLastScrapedTime := m.C.Now().Add(-24 * time.Hour)
	lastScrape, err := m.Queries.GetLastScrape(ctx, sourceID)

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return defaultLastScrapedTime, errors.Wrap(err, "problem getting last scrape")
		}

		return defaultLastScrapedTime, nil
	}

	return lastScrape.CreatedAt, nil
}
