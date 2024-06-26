package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/models"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

type SpyScraper struct {
	calls int
}

func NewSpyScraper() *SpyScraper {
	return &SpyScraper{
		calls: 0,
	}
}

func (m *SpyScraper) ScrapeJobs(lastScraped time.Time) ([]models.Job, []error) {
	m.calls++
	return []models.Job{}, []error{}
}

func (m *SpyScraper) GetSourceID() string {
	return "linkedin"
}

type SpyFailingScraper struct {
	*SpyScraper
}

func (m *SpyFailingScraper) ScrapeJobs(lastScraped time.Time) ([]models.Job, []error) {
	m.calls++
	return []models.Job{}, []error{fmt.Errorf("could not scrape heading"), fmt.Errorf("problem scraping paragraph")}
}

func (m *SpyFailingScraper) GetSourceID() string {
	return "linkedin"
}

func newSpyFailingScraper() *SpyFailingScraper {
	ms := NewSpyScraper()
	return &SpyFailingScraper{
		SpyScraper: ms,
	}
}

const callsBetween8pmAnd10pm = 5
const callsBetween7amAnd8am = 3

func TestScrapeService(t *testing.T) {
	t.Run("calls the scrapers the correct amount of times given a cron", func(t *testing.T) {
		c := createFakeClock(t, "America/Denver", func(loc *time.Location) time.Time {
			return time.Date(2024, time.May, 19, 20, 0, 0, 0, loc)
		})
		logBufferSpy := new(bytes.Buffer)
		logger := slog.New(slog.NewTextHandler(logBufferSpy, nil))

		db, _ := database.NewInMemoryDB()

		scrapeService := newTestScrapeService(logger, db, c)

		scraper1 := NewSpyScraper()
		scraper2 := NewSpyScraper()
		scrapers := []models.ScraperModelInterface{scraper1, scraper2}
		scrapeService.addScrapers(scrapers)
		go scrapeService.start("TZ=America/Denver */30 7-22 * * *", false)
		c.BlockUntil(1)

		// loop one extra time to ensure no extra calls are made
		callNotMade := simulateCron(c, callsBetween8pmAnd10pm+1, 30*time.Minute)
		assert.Equal(t, callsBetween8pmAnd10pm, scraper1.calls)
		assertCallNotMade(t, callNotMade, scraper1.calls)
		assert.Equal(t, callsBetween8pmAnd10pm, scraper2.calls)
		assertCallNotMade(t, callNotMade, scraper2.calls)

		// advance to 6:30am
		c.Advance(7*time.Hour + 30*time.Minute)

		callNotMade = simulateCron(c, callsBetween7amAnd8am+1, 30*time.Minute)
		assert.Equal(t, callsBetween8pmAnd10pm+callsBetween7amAnd8am, scraper1.calls)
		assertCallNotMade(t, callNotMade, scraper1.calls)
		assert.Equal(t, callsBetween8pmAnd10pm+callsBetween7amAnd8am, scraper2.calls)
		assertCallNotMade(t, callNotMade, scraper2.calls)
	})

	t.Run("gets the latest scrape", func(t *testing.T) {
		c := createFakeClock(t, "America/Denver", func(loc *time.Location) time.Time {
			return time.Date(2024, time.May, 19, 20, 0, 0, 0, loc)
		})

		ctx := context.Background()
		logBufferSpy := new(bytes.Buffer)
		logger := slog.New(slog.NewTextHandler(logBufferSpy, nil))

		db, _ := database.NewInMemoryDB()

		queries := database.New(db)
		scrapes := &models.ScrapeModel{Queries: queries, C: c}

		scrapeService := newTestScrapeService(logger, db, c)
		scraper1 := NewSpyScraper()
		scrapers := []models.ScraperModelInterface{scraper1}
		scrapeService.addScrapers(scrapers)

		go scrapeService.start("TZ=America/Denver */30 7-22 * * *", false)
		c.BlockUntil(1)

		simulateCron(c, callsBetween8pmAnd10pm+1, 30*time.Minute)

		scrape, err := scrapes.Latest(ctx, "linkedin")
		assert.NoError(t, err)
		assert.Equal(t, callsBetween8pmAnd10pm, scrape.ID)
	})

	t.Run("logs errors that occur", func(t *testing.T) {
		c := createFakeClock(t, "America/Denver", func(loc *time.Location) time.Time {
			return time.Date(2024, time.May, 19, 20, 0, 0, 0, loc)
		})
		logBufferSpy := new(bytes.Buffer)
		logger := slog.New(slog.NewTextHandler(logBufferSpy, nil))

		db, _ := database.NewInMemoryDB()

		scrapeService := newTestScrapeService(logger, db, c)

		scraper := newSpyFailingScraper()
		scrapers := []models.ScraperModelInterface{scraper}
		scrapeService.addScrapers(scrapers)

		go scrapeService.start("TZ=America/Denver */30 7-22 * * *", false)
		c.BlockUntil(1)

		simulateCron(c, 2, 30*time.Minute)

		assert.Contains(t, logBufferSpy.String(), "could not scrape heading")
		assert.Contains(t, logBufferSpy.String(), "problem scraping paragraph")
	})

	t.Run("calls scraper immediately when specified", func(t *testing.T) {
		c := createFakeClock(t, "America/Denver", func(loc *time.Location) time.Time {
			return time.Date(2024, time.May, 19, 7, 10, 0, 0, loc)
		})
		logBufferSpy := new(bytes.Buffer)
		logger := slog.New(slog.NewTextHandler(logBufferSpy, nil))

		db, _ := database.NewInMemoryDB()

		scrapeService := newTestScrapeService(logger, db, c)

		scraper1 := NewSpyScraper()
		scrapers := []models.ScraperModelInterface{scraper1}
		scrapeService.addScrapers(scrapers)

		go scrapeService.start("TZ=America/Denver */30 7-22 * * *", true)
		c.BlockUntil(1)

		simulateCron(c, 3, 30*time.Minute)
		assert.Equal(t, 3, scraper1.calls)
	})
}

func createFakeClock(t *testing.T, location string, getTime func(*time.Location) time.Time) clockwork.FakeClock {
	t.Helper()
	loadedLocation, err := time.LoadLocation(location)

	if err != nil {
		t.Fatal("failed to load location")
	}

	convertedStartTime := getTime(loadedLocation)
	c := clockwork.NewFakeClockAt(convertedStartTime)

	return c
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

func assertCallNotMade(t *testing.T, callNotMade bool, calls int) {
	t.Helper()
	if !assert.Equal(t, true, callNotMade, fmt.Sprintf("should not make calls outside of cron interval\nscraper.calls is %v, expected %v", calls, callsBetween8pmAnd10pm+callsBetween7amAnd8am)) {
		t.FailNow()
	}
}

func newTestScrapeService(logger *slog.Logger, db *sql.DB, c clockwork.Clock) *scrapeService {
	queries := database.New(db)

	companies := &models.CompanyModel{Queries: queries}
	jobs := &models.JobModel{Queries: queries, Companies: companies}
	scrapes := &models.ScrapeModel{Queries: queries, C: c}

	httpClient := newHttpClient(logger, &config{})

	return &scrapeService{logger: logger, httpClient: httpClient, scrapes: scrapes, jobs: jobs, c: c}
}
