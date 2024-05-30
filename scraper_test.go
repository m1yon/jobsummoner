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
		location, err := time.LoadLocation("America/Denver")
		startTime := time.Date(2024, time.May, 19, 20, 0, 0, 0, location)

		if err != nil {
			t.Fatal("failed to load location")
		}

		c := clockwork.NewFakeClockAt(startTime)

		scraper := NewMockScraper()

		go ScrapeLoop(c, scraper, "TZ=America/Denver */30 7-22 * * *")

		c.BlockUntil(1)

		callNotMade := false
		// loop one extra time to ensure no extra calls are made
		for i := 0; i < callsBetween8pmAnd10pm+1; i++ {
			c.Advance(30 * time.Minute)

			select {
			case <-scraper.ch:
				if !assert.Equal(t, i+1, scraper.calls) {
					t.FailNow()
				}
			case <-time.After(50 * time.Millisecond):
				callNotMade = true
			}
		}

		assert.Equal(t, callsBetween8pmAnd10pm, scraper.calls)
		if !assert.Equal(t, true, callNotMade, fmt.Sprintf("should stop making calls when the end of the cron interval is reached\nscraper.calls is %v, expected %v", scraper.calls, callsBetween8pmAnd10pm)) {
			t.FailNow()
		}

		// advance to 6:30am
		c.Advance(7*time.Hour + 30*time.Minute)

		callNotMade = false
		for i := 0; i < callsBetween7amAnd8am+1; i++ {
			fmt.Println(c.Now())

			select {
			case <-scraper.ch:
				if i == 0 {
					assert.FailNow(t, "cron triggered at 6:30am")
				}

				if !assert.Equal(t, callsBetween8pmAnd10pm+i, scraper.calls) {
					t.FailNow()
				}
			case <-time.After(50 * time.Millisecond):
				callNotMade = true
			}

			c.Advance(30 * time.Minute)
		}

		assert.Equal(t, callsBetween8pmAnd10pm+callsBetween7amAnd8am, scraper.calls)
		if !assert.Equal(t, true, callNotMade, fmt.Sprintf("should stop making calls when the end of the cron interval is reached\nscraper.calls is %v, expected %v", scraper.calls, callsBetween8pmAnd10pm+callsBetween7amAnd8am)) {
			t.FailNow()
		}
	})
}
