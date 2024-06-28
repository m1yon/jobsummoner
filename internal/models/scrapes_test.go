package models

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/database"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestScrapes(t *testing.T) {
	t.Run("can create scrape and immediately get scrape", func(t *testing.T) {
		ctx := context.Background()
		scrapes, c := newTestScrapeModel()

		err := scrapes.Create(ctx, "linkedin", c.Now())
		assert.NoError(t, err)

		scrape, err := scrapes.Latest(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, "linkedin", scrape.SourceID)
		}
	})

	t.Run("gets latest scrape", func(t *testing.T) {
		ctx := context.Background()
		scrapes, c := newTestScrapeModel()

		_ = scrapes.Create(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapes.Create(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapes.Create(ctx, "linkedin", c.Now())

		scrape, err := scrapes.Latest(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, "linkedin", scrape.SourceID)
			assert.Equal(t, 3, scrape.ID)
		}
	})

	t.Run("gets latest scrape time", func(t *testing.T) {
		ctx := context.Background()
		scrapes, c := newTestScrapeModel()

		_ = scrapes.Create(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapes.Create(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapes.Create(ctx, "linkedin", c.Now())

		scrapeTime, err := scrapes.LastRan(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, c.Now().UTC(), scrapeTime)
		}
	})
}

func newTestScrapeModel() (*ScrapeModel, clockwork.FakeClock) {
	c := clockwork.NewFakeClock()

	db, _ := database.NewInMemoryDB()
	queries := database.New(db)

	scrapes := &ScrapeModel{Queries: queries, C: c}

	return scrapes, c
}
