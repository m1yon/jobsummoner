package jobsummoner

import (
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

type MockScraper struct {
	ch    chan interface{}
	calls int
}

func (m *MockScraper) ScrapeJobs() (ScrapedJobsResults, []error) {
	m.calls++
	m.ch <- struct{}{}

	return ScrapedJobsResults{}, []error{}
}

func NewMockScraper() *MockScraper {
	ch := make(chan interface{})

	return &MockScraper{
		ch:    ch,
		calls: 0,
	}
}

const callsBetween8pmAnd10pm = 5
const callsBetween7amAnd8am = 3

func TestScrapeLoop(t *testing.T) {
	t.Run("calls the function correctly on a cron", func(t *testing.T) {
		c := getFakeClock(t)
		scraper := NewMockScraper()

		go ScrapeLoop(c, scraper, "TZ=America/Denver */30 7-22 * * *")
		c.BlockUntil(1)

		// loop one extra time to ensure no extra calls are made
		callNotMade := simulateCron(c, scraper, callsBetween8pmAnd10pm+1, 30*time.Minute)
		assert.Equal(t, callsBetween8pmAnd10pm, scraper.calls)
		assertCallNotMade(t, callNotMade, scraper.calls)

		// advance to 6:30am
		c.Advance(7*time.Hour + 30*time.Minute)

		callNotMade = simulateCron(c, scraper, callsBetween7amAnd8am+1, 30*time.Minute)
		assert.Equal(t, callsBetween8pmAnd10pm+callsBetween7amAnd8am, scraper.calls)
		assertCallNotMade(t, callNotMade, scraper.calls)
	})
}

func simulateCron(c clockwork.FakeClock, scraper *MockScraper, numberOfCalls int, advanceInterval time.Duration) bool {
	missed := false
	for i := 0; i < numberOfCalls; i++ {
		select {
		case <-scraper.ch:
		case <-time.After(50 * time.Millisecond):
			missed = true
		}

		c.Advance(advanceInterval)
	}

	return missed
}

func getFakeClock(t *testing.T) clockwork.FakeClock {
	t.Helper()
	location, err := time.LoadLocation("America/Denver")
	startTime := time.Date(2024, time.May, 19, 20, 0, 0, 0, location)

	if err != nil {
		t.Fatal("failed to load location")
	}

	c := clockwork.NewFakeClockAt(startTime)
	return c
}

func assertCallNotMade(t *testing.T, callNotMade bool, calls int) {
	t.Helper()
	if !assert.Equal(t, true, callNotMade, fmt.Sprintf("should not make calls outside of cron interval\nscraper.calls is %v, expected %v", calls, callsBetween8pmAnd10pm+callsBetween7amAnd8am)) {
		t.FailNow()
	}
}
