package jobsummoner

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

type MockScraper struct {
	ch    chan bool
	calls int
}

func (m *MockScraper) ScrapeJobs() (ScrapedJobsResults, []error) {
	newCallsValue := m.calls + 1
	m.calls = newCallsValue

	m.ch <- true

	return ScrapedJobsResults{}, []error{}
}

func NewMockScraper() *MockScraper {
	ch := make(chan bool)

	return &MockScraper{
		ch:    ch,
		calls: 0,
	}
}

const callsBetween5pmAnd10pm = 11

func TestScrapeLoop(t *testing.T) {
	location, err := time.LoadLocation("America/Denver")
	startTime := time.Date(2024, time.May, 19, 17, 0, 0, 0, location)

	if err != nil {
		t.Fatal("failed to load location")
	}

	c := clockwork.NewFakeClockAt(startTime)

	scraper := NewMockScraper()

	go ScrapeLoop(c, scraper, "TZ=America/Denver */30 7-22 * * *")

	c.BlockUntil(1)

	// loop one extra time to ensure no extra calls are made
	for i := 0; i < callsBetween5pmAnd10pm+1; i++ {
		c.Advance(30 * time.Minute)

		select {
		case <-scraper.ch:
			assert.Equal(t, i+1, scraper.calls)
		case <-time.After(50 * time.Millisecond):
			break
		}
	}

	assert.Equal(t, callsBetween5pmAnd10pm, scraper.calls)
}
