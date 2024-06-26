package sqlitedb

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestScrapeRepository(t *testing.T) {
	t.Run("can create scrape and immediately get scrape", func(t *testing.T) {
		ctx := context.Background()
		db, _ := NewInMemoryDB()
		c := clockwork.NewFakeClock()
		scrapeRepository := NewSqliteScrapeRepository(db, c)

		err := scrapeRepository.CreateScrape(ctx, "linkedin", c.Now())
		assert.NoError(t, err)

		scrape, err := scrapeRepository.GetLastScrape(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, "linkedin", scrape.SourceID)
		}
	})

	t.Run("gets latest scrape", func(t *testing.T) {
		ctx := context.Background()
		db, _ := NewInMemoryDB()
		c := clockwork.NewFakeClock()
		scrapeRepository := NewSqliteScrapeRepository(db, c)

		_ = scrapeRepository.CreateScrape(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapeRepository.CreateScrape(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapeRepository.CreateScrape(ctx, "linkedin", c.Now())

		scrape, err := scrapeRepository.GetLastScrape(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, "linkedin", scrape.SourceID)
			assert.Equal(t, 3, scrape.ID)
		}
	})

	t.Run("gets latest scrape time", func(t *testing.T) {
		ctx := context.Background()
		db, _ := NewInMemoryDB()
		c := clockwork.NewFakeClock()
		scrapeRepository := NewSqliteScrapeRepository(db, c)

		_ = scrapeRepository.CreateScrape(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapeRepository.CreateScrape(ctx, "linkedin", c.Now())
		c.Advance(time.Hour * 1)
		_ = scrapeRepository.CreateScrape(ctx, "linkedin", c.Now())

		scrapeTime, err := scrapeRepository.GetLastScrapeTime(ctx, "linkedin")
		if assert.NoError(t, err) {
			assert.Equal(t, c.Now().UTC(), scrapeTime)
		}
	})
}
