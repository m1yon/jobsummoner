package jobsummoner

import (
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

type MockScraper struct {
	calls int
}

func (m *MockScraper) ScrapeJobs() (ScrapedJobsResults, []error) {
	m.calls++
	return ScrapedJobsResults{}, []error{}
}

func NewMockScraper() *MockScraper {
	return &MockScraper{
		calls: 0,
	}
}

func TestScrapeLoop(t *testing.T) {
	c := clockwork.NewFakeClock()
	interval := 5 * time.Minute

	scraper := NewMockScraper()

	go ScrapeLoop(c, scraper, interval)

	fmt.Printf("now %v\n", c.Now())

	c.BlockUntil(1)

	assert.Equal(t, 0, scraper.calls)

	c.Advance(interval + time.Second)
	c.BlockUntil(1)
	assert.Equal(t, 1, scraper.calls)

	c.Advance(interval)
	c.BlockUntil(1)
	assert.Equal(t, 2, scraper.calls)
}
