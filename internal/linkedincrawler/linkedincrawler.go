package linkedincrawler

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/stealth"
	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/robfig/cron"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func ScrapeLoop(db *sql.DB) {
	ctx := context.Background()
	dbQueries := database.New(db)

	scrapes, err := dbQueries.GetAllScrapesWithKeywords(ctx)

	if err != nil {
		if !strings.Contains(err.Error(), "converting NULL to int64 is unsupported") {
			slog.Error("failed querying for all scrapes", tint.Err(err))
			return
		}

		slog.Warn("no scrapes found in db")
	}

	location, err := time.LoadLocation("America/Denver")

	if err != nil {
		slog.Error("failed creating cron location", tint.Err(err))
	}

	c := cron.NewWithLocation(location)

	for _, currentScrape := range scrapes {
		c.AddFunc("0 */30 7-22 * * *", func() {
			localOptions := currentScrape
			scrape(db, localOptions)
		})
	}

	scrapeOnStart := os.Getenv("SCRAPE_ON_START")

	if scrapeOnStart == "true" {
		for _, currentScrape := range scrapes {
			scrape(db, currentScrape)
		}
	}

	c.Start()
}

var (
	postingsScrapedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "jobsummoner_postings_scraped_total",
		Help: "The total number of job postings scraped.",
	})
	newPostingsScrapedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "jobsummoner_new_postings_scraped_total",
		Help: "The total number of new job postings scraped.",
	})
)

func scrape(db *sql.DB, currentScrape database.GetAllScrapesWithKeywordsRow) {
	ctx := context.Background()
	dbQueries := database.New(db)

	PROXY_HOSTNAME := os.Getenv("PROXY_HOSTNAME")
	PROXY_USERNAME := os.Getenv("PROXY_USERNAME")
	PROXY_PASSWORD := os.Getenv("PROXY_PASSWORD")

	proxyEnabled := len(PROXY_HOSTNAME) != 0

	slog.Info("starting scrape", slog.String("name", currentScrape.Name), slog.Bool("proxy", proxyEnabled))

	defer func() {
		if err := dbQueries.UpdateLastScraped(ctx, currentScrape.ID); err != nil {
			slog.Error("failed to updated last_scraped field", tint.Err(err))
		}
	}()

	// proxy setup
	l := launcher.New()

	if proxyEnabled {
		l = l.Set(flags.ProxyServer, PROXY_HOSTNAME)
	}
	controlURL, _ := l.Launch()
	browser := rod.New()
	err := browser.ControlURL(controlURL).Connect()

	if err != nil {
		slog.Error("browser connection failed", tint.Err(err))
		return
	}

	defer browser.MustClose()

	if proxyEnabled {
		go browser.MustHandleAuth(PROXY_USERNAME, PROXY_PASSWORD)()
	}

	page, err := stealth.Page(browser)

	if err != nil {
		slog.Error("failed to create stealth page", tint.Err(err))
		return
	}

	url, err := url.Parse("https://linkedin.com/jobs/search/")

	if err != nil {
		slog.Error("failed to build URL", tint.Err(err))
		return
	}

	defaultTime := time.Hour * 1

	q := url.Query()
	q.Set("keywords", currentScrape.Keywords)
	q.Set("location", currentScrape.Location)
	q.Set("f_WT", fmt.Sprintf("%v", currentScrape.WorkType))
	q.Set("f_JT", "F")
	q.Set("f_SB2", "5")
	q.Set("f_TPR", fmt.Sprintf("r%v", defaultTime.Seconds()))

	url.RawQuery = q.Encode()

	slog.Debug(url.String())

	err = page.Navigate(url.String())

	if err != nil {
		slog.Error("failed to navigate to linkedin", slog.String("url", url.String()), tint.Err(err))
		return
	}

	page.MustWaitStable()

	// scroll to the bottom of the page to load all virtualized resources
	scrollHeight := page.MustEval("() => document.documentElement.scrollHeight").Int()
	page.Mouse.Scroll(0.0, float64(scrollHeight), 20)

	jobPostings, err := page.Elements(".jobs-search__results-list > li")

	if err != nil {
		slog.Error("failed to query for job postings", slog.String("url", url.String()), tint.Err(err))
		return
	}

	numberOfJobPostings := len(jobPostings)
	numberOfJobRepostings := 0

	for _, jobPosting := range jobPostings {
		position, err := jobPosting.Element(".base-search-card__title")

		if err != nil {
			slog.Error("failed to query for position in job posting", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		positionText, err := position.Text()

		if err != nil {
			slog.Error("failed to get position text from element", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		companyName, err := jobPosting.Element(".base-search-card__subtitle")

		if err != nil {
			slog.Error("failed to query for company name in job posting", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		companyNameText, err := companyName.Text()

		if err != nil {
			slog.Error("failed to get company name text from element", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		location, err := jobPosting.Element(".job-search-card__location")

		if err != nil {
			slog.Error("failed to query for location in job posting", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		locationText, err := location.Text()

		if err != nil {
			slog.Error("failed to get location text from element", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		postingURL, err := jobPosting.Element("a")

		if err != nil {
			slog.Error("failed to query for posting url in job posting", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		postingUrlText, err := postingURL.Property("href")

		if err != nil {
			slog.Error("failed to get url from element", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		relativeListingDate, err := jobPosting.Element(".base-search-card__metadata > time")

		if err != nil {
			slog.Error("faled to query for relative listing date in job posting", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		relativeListingDateText, err := relativeListingDate.Text()

		if err != nil {
			slog.Error("failed to get url from element", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		listingDate, err := parseRelativeTime(relativeListingDateText)

		if err != nil {
			slog.Error("failed to parse relative listing date", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		companyLink, err := jobPosting.Element(".base-search-card__subtitle > a")

		if err != nil {
			slog.Error("failed to query for company link in job posting", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		companyLinkURL, err := companyLink.Property("href")

		if err != nil {
			slog.Error("failed to get company url from element", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		parsedCompanyLinkURL, err := url.Parse(companyLinkURL.String())

		if err != nil {
			slog.Error("failed parsing company link url", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		segments := strings.Split(parsedCompanyLinkURL.EscapedPath(), "/")
		companySlug := segments[len(segments)-1]

		companyAvatar, err := jobPosting.Element(".base-card img")

		if err != nil {
			slog.Error("failed to query for company avatar in job posting", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		companyAvatarSrc, err := companyAvatar.Property("src")

		if err != nil {
			slog.Error("failed parsing company avatar", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		err = dbQueries.CreateCompany(ctx, database.CreateCompanyParams{ID: companySlug, Name: companyNameText, Url: companyLinkURL.String(), Avatar: sql.NullString{String: companyAvatarSrc.String(), Valid: true}})

		if err != nil {
			slog.Error("failed inserting company", slog.String("url", url.String()), tint.Err(err))
			continue
		}

		jobPostingID := getJobPostingID(companySlug, positionText)

		err = dbQueries.CreateJobPosting(ctx, database.CreateJobPostingParams{
			ID:           jobPostingID,
			Position:     positionText,
			Url:          postingUrlText.String(),
			CompanyID:    companySlug,
			LastPosted:   listingDate.UTC(),
			Location:     sql.NullString{String: locationText, Valid: true},
			LocationType: currentScrape.WorkType,
		})

		if err != nil {
			if sqlError, ok := err.(*sqlite.Error); ok {
				// fail if it's not an expected error
				if sqlError.Code() != sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
					slog.Error("failed inserting job posting", slog.String("url", url.String()), tint.Err(err))
					continue
				}

				// if the posting already exists, just update the last_posted field
				numberOfJobRepostings++

				err = dbQueries.UpdateJobPostingLastPosted(ctx, database.UpdateJobPostingLastPostedParams{
					ID:         jobPostingID,
					LastPosted: listingDate.UTC(),
				})

				if err != nil {
					slog.Error("failed updating job posting's last_posted field", slog.String("url", url.String()), tint.Err(err))
					continue
				}
			}
		}

		err = dbQueries.CreateUserJobPosting(ctx, database.CreateUserJobPostingParams{
			JobPostingID: jobPostingID,
			UserID:       int64(currentScrape.UserID),
		})

		if err != nil {
			if sqlError, ok := err.(*sqlite.Error); ok {
				if sqlError.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {

					postingsScrapedTotal.Inc()
					continue
				}
			}

			slog.Error("failed insert user job posting", slog.String("url", url.String()), tint.Err(err))
		}

		postingsScrapedTotal.Inc()
		newPostingsScrapedTotal.Inc()
	}

	slog.Info("scrape finished", slog.String("name", currentScrape.Name), slog.Int("postings", numberOfJobPostings), slog.Int("repostings", numberOfJobRepostings))
}

func getJobPostingID(company_id string, position string) string {
	data := company_id + "|" + position

	hasher := sha256.New()

	hasher.Write([]byte(data))

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}

func parseRelativeTime(input string) (time.Time, error) {
	// Regex to find number and unit
	re := regexp.MustCompile(`(\d+)\s*(second|minute|hour|day)s?\s*ago`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 3 {
		return time.Time{}, fmt.Errorf("invalid format")
	}

	// Convert number to integer
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, err
	}

	// Determine the unit of time
	var duration time.Duration
	switch matches[2] {
	case "second":
		duration = time.Second * time.Duration(value)
	case "minute":
		duration = time.Minute * time.Duration(value)
	case "hour":
		duration = time.Hour * time.Duration(value)
	case "day":
		duration = time.Hour * 24 * time.Duration(value)
	default:
		return time.Time{}, fmt.Errorf("unknown time unit")
	}

	// Calculate the time
	actualTime := time.Now().Add(-duration)
	return actualTime, nil
}
