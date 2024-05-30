package scrapeservice

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

type MockScraper struct {
	calls int
}

func (m *MockScraper) ScrapeJobs() (jobsummoner.ScrapedJobsResults, []error) {
	m.calls++
	return jobsummoner.ScrapedJobsResults{}, []error{}
}

func NewMockScraper() *MockScraper {
	return &MockScraper{
		calls: 0,
	}
}

type MockFailingScraper struct {
	*MockScraper
}

func (m *MockFailingScraper) ScrapeJobs() (jobsummoner.ScrapedJobsResults, []error) {
	m.calls++
	return jobsummoner.ScrapedJobsResults{}, []error{fmt.Errorf("could not scrape heading"), fmt.Errorf("problem scraping paragraph")}
}

func newMockFailingScraper() *MockFailingScraper {
	ms := NewMockScraper()
	return &MockFailingScraper{
		MockScraper: ms,
	}
}

const callsBetween8pmAnd10pm = 5
const callsBetween7amAnd8am = 3

func TestScrapeLoop(t *testing.T) {
	t.Run("calls the function correctly on a cron", func(t *testing.T) {
		c := getFakeClock(t)
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		scrapeService := NewDefaultScrapeService(c, logger)

		scraper := NewMockScraper()
		go scrapeService.Start(scraper, "TZ=America/Denver */30 7-22 * * *")
		c.BlockUntil(1)

		// loop one extra time to ensure no extra calls are made
		callNotMade := simulateCron(c, callsBetween8pmAnd10pm+1, 30*time.Minute)
		assert.Equal(t, callsBetween8pmAnd10pm, scraper.calls)
		assertCallNotMade(t, callNotMade, scraper.calls)

		// advance to 6:30am
		c.Advance(7*time.Hour + 30*time.Minute)

		callNotMade = simulateCron(c, callsBetween7amAnd8am+1, 30*time.Minute)
		assert.Equal(t, callsBetween8pmAnd10pm+callsBetween7amAnd8am, scraper.calls)
		assertCallNotMade(t, callNotMade, scraper.calls)
	})

	t.Run("logs errors that occur", func(t *testing.T) {
		c := getFakeClock(t)
		logBufferSpy := new(bytes.Buffer)
		logger := slog.New(slog.NewTextHandler(logBufferSpy, nil))
		scrapeService := NewDefaultScrapeService(c, logger)

		scraper := newMockFailingScraper()

		go scrapeService.Start(scraper, "TZ=America/Denver */30 7-22 * * *")
		c.BlockUntil(1)

		simulateCron(c, 2, 30*time.Minute)

		assert.Contains(t, logBufferSpy.String(), "could not scrape heading")
		assert.Contains(t, logBufferSpy.String(), "problem scraping paragraph")
	})
}

func simulateCron(c clockwork.FakeClock, numberOfCalls int, advanceInterval time.Duration) bool {
	ch := make(chan interface{})

	missed := false
	for i := 0; i < numberOfCalls; i++ {
		select {
		case <-ch:
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
