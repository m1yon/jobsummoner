package models

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestScrapeRepository(t *testing.T) {
	t.Run("can create scrape and immediately get scrape", func(t *testing.T) {
		ctx := context.Background()
		c := clockwork.NewFakeClock()

		db, _ := sqlitedb.NewInMemoryDB()
		queries := sqlitedb.New(db)

		scrapes := &ScrapeModel{Queries: queries, C: c}

		err := scrapes.CreateScrape(ctx, "linkedin", c.Now())
		assert.NoError(t, err)

		scrape, err := scrapes.GetLastScrape(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, "linkedin", scrape.SourceID)
		}
	})

	t.Run("gets latest scrape", func(t *testing.T) {
		ctx := context.Background()
		c := clockwork.NewFakeClock()

		db, _ := sqlitedb.NewInMemoryDB()
		queries := sqlitedb.New(db)

		scrapes := &ScrapeModel{Queries: queries, C: c}

		_ = scrapes.CreateScrape(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapes.CreateScrape(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapes.CreateScrape(ctx, "linkedin", c.Now())

		scrape, err := scrapes.GetLastScrape(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, "linkedin", scrape.SourceID)
			assert.Equal(t, 3, scrape.ID)
		}
	})

	t.Run("gets latest scrape time", func(t *testing.T) {
		ctx := context.Background()
		c := clockwork.NewFakeClock()

		db, _ := sqlitedb.NewInMemoryDB()
		queries := sqlitedb.New(db)

		scrapes := &ScrapeModel{Queries: queries, C: c}

		_ = scrapes.CreateScrape(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapes.CreateScrape(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapes.CreateScrape(ctx, "linkedin", c.Now())

		scrapeTime, err := scrapes.GetLastScrapeTime(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, c.Now().UTC(), scrapeTime)
		}
	})
}
