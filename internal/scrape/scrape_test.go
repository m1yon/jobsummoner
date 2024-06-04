package scrape

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SpyScraper struct {
	calls int
}

func (m *SpyScraper) ScrapeJobs() ([]jobsummoner.Job, []error) {
	m.calls++
	return []jobsummoner.Job{}, []error{}
}

func NewSpyScraper() *SpyScraper {
	return &SpyScraper{
		calls: 0,
	}
}

type SpyFailingScraper struct {
	*SpyScraper
}

func (m *SpyFailingScraper) ScrapeJobs() ([]jobsummoner.Job, []error) {
	m.calls++
	return []jobsummoner.Job{}, []error{fmt.Errorf("could not scrape heading"), fmt.Errorf("problem scraping paragraph")}
}

func newSpyFailingScraper() *SpyFailingScraper {
	ms := NewSpyScraper()
	return &SpyFailingScraper{
		SpyScraper: ms,
	}
}

const callsBetween8pmAnd10pm = 5
const callsBetween7amAnd8am = 3

type mockJobService struct {
	mock.Mock
}

func (j *mockJobService) GetJob(ctx context.Context, id string) (jobsummoner.Job, error) {
	j.Called()
	return jobsummoner.Job{}, nil
}

func (j *mockJobService) GetJobs(ctx context.Context) ([]jobsummoner.Job, []error) {
	j.Called()
	return []jobsummoner.Job{}, nil
}

func (j *mockJobService) CreateJobs(ctx context.Context, jobs []jobsummoner.Job) ([]string, []error) {
	j.Called()
	return nil, nil
}

func (j *mockJobService) CreateJob(ctx context.Context, jobs jobsummoner.Job) (string, error) {
	j.Called()
	return "", nil
}

func TestScrapeService(t *testing.T) {
	t.Run("calls multiple scrapers correctly on a cron and sends the results to the Job Service", func(t *testing.T) {
		c, _, jobServiceMock, scrapeService := initScrapeServiceSpies(t)
		scraper1 := NewSpyScraper()
		scraper2 := NewSpyScraper()
		scrapers := []jobsummoner.Scraper{scraper1, scraper2}

		go scrapeService.Start(scrapers, "TZ=America/Denver */30 7-22 * * *")
		c.BlockUntil(1)

		// loop one extra time to ensure no extra calls are made
		callNotMade := simulateCron(c, callsBetween8pmAnd10pm+1, 30*time.Minute)
		// expect an additional startup call
		assert.Equal(t, callsBetween8pmAnd10pm+1, scraper1.calls)
		assertCallNotMade(t, callNotMade, scraper1.calls)
		assert.Equal(t, callsBetween8pmAnd10pm+1, scraper2.calls)
		assertCallNotMade(t, callNotMade, scraper2.calls)

		// advance to 6:30am
		c.Advance(7*time.Hour + 30*time.Minute)

		callNotMade = simulateCron(c, callsBetween7amAnd8am+1, 30*time.Minute)
		// expect an additional startup call
		assert.Equal(t, callsBetween8pmAnd10pm+callsBetween7amAnd8am+1, scraper1.calls)
		assertCallNotMade(t, callNotMade, scraper1.calls)
		assert.Equal(t, callsBetween8pmAnd10pm+callsBetween7amAnd8am+1, scraper2.calls)
		assertCallNotMade(t, callNotMade, scraper2.calls)

		jobServiceMock.AssertExpectations(t)
		assert.Equal(t, scraper1.calls+scraper2.calls, len(jobServiceMock.Calls))
	})

	t.Run("logs errors that occur", func(t *testing.T) {
		c, logBufferSpy, jobServiceMock, scrapeService := initScrapeServiceSpies(t)
		scraper := newSpyFailingScraper()
		scrapers := []jobsummoner.Scraper{scraper}

		go scrapeService.Start(scrapers, "TZ=America/Denver */30 7-22 * * *")
		c.BlockUntil(1)

		simulateCron(c, 2, 30*time.Minute)

		assert.Contains(t, logBufferSpy.String(), "could not scrape heading")
		assert.Contains(t, logBufferSpy.String(), "problem scraping paragraph")

		jobServiceMock.AssertExpectations(t)
		assert.Equal(t, scraper.calls, len(jobServiceMock.Calls))
	})
}

func initScrapeServiceSpies(t *testing.T) (clockwork.FakeClock, *bytes.Buffer, *mockJobService, *DefaultScrapeService) {
	t.Helper()
	c := getFakeClock(t)
	logBufferSpy := new(bytes.Buffer)
	logger := slog.New(slog.NewTextHandler(logBufferSpy, nil))
	jobServiceMock := new(mockJobService)
	scrapeService := NewDefaultScrapeService(c, logger, jobServiceMock)

	jobServiceMock.On("CreateJobs", mock.Anything).Return()

	return c, logBufferSpy, jobServiceMock, scrapeService
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
