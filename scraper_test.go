package jobsummoner

import (
	"runtime"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
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
	c := clock.NewMock()
	interval := 5 * time.Minute
	scraper := NewMockScraper()

	go ScrapeLoop(c, scraper, interval)

	runtime.Gosched()

	assert.Equal(t, scraper.calls, 0)

	c.Add(interval)
	assert.Equal(t, scraper.calls, 1)

	c.Add(interval)
	assert.Equal(t, scraper.calls, 2)
}
