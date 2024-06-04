package sqlitedb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScrapeRepository(t *testing.T) {
	t.Run("can create scrape and immediately get scrape", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		scrapeRepository := NewSqliteScrapeRepository(db)

		err := scrapeRepository.CreateScrape(ctx, "linkedin", time.Now())
		assert.NoError(t, err)

		scrape, err := scrapeRepository.GetLastScrape(ctx, "linkedin")
		assert.NoError(t, err)
		assert.Equal(t, "linkedin", scrape.SourceID)
	})

	t.Run("gets latest scrape", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		scrapeRepository := NewSqliteScrapeRepository(db)

		err := scrapeRepository.CreateScrape(ctx, "linkedin", time.Now())
		assert.NoError(t, err)
		err = scrapeRepository.CreateScrape(ctx, "linkedin", time.Now().Add(time.Hour*1))
		assert.NoError(t, err)
		err = scrapeRepository.CreateScrape(ctx, "linkedin", time.Now().Add(time.Hour*2))
		assert.NoError(t, err)

		scrape, err := scrapeRepository.GetLastScrape(ctx, "linkedin")
		assert.NoError(t, err)
		assert.Equal(t, "linkedin", scrape.SourceID)
		assert.Equal(t, 3, scrape.ID)
	})
}
