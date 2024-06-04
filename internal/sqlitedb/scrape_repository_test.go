package sqlitedb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrapeRepository(t *testing.T) {
	t.Run("can create scrape", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		scrapeRepository := NewSqliteScrapeRepository(db)

		err := scrapeRepository.CreateScrape(ctx, "linkedin")
		assert.NoError(t, err)
	})
}
