package linkedincrawler

import (
	"context"
	"database/sql"
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
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func ScrapeLoop(db *sql.DB) {
	scrapes := []scrapeOptions{
		{
			name:         "Remote Roles",
			keywords:     []string{"typescript", "react"},
			location:     "United States",
			workType:     2,             // 2 = remote
			jobTypes:     []string{"F"}, // F = fulltime
			salaryRanges: []string{"5"}, // 5 = $120,000+
			ageOfPosting: 24 * time.Hour,
		},
		{
			name:         "Colorado Hybrid Roles",
			keywords:     []string{"typescript", "react"},
			location:     "Colorado, United States",
			workType:     3,
			jobTypes:     []string{"F"},
			salaryRanges: []string{"5"},
			ageOfPosting: 24 * time.Hour,
		},
	}

	for {
		for _, options := range scrapes {
			slog.Info(fmt.Sprintf("starting scrape for '%v'...", options.name))

			numberOfJobPostings, numberOfJobRepostings, err := scrape(db, options)

			if err != nil {
				slog.Error(fmt.Sprintf("failed to scrape '%v'", options.name), tint.Err(err))
			}

			slog.Info(fmt.Sprintf("successfully scraped '%v' and got (%v) job postings, (%v) of them being reposts", options.name, numberOfJobPostings, numberOfJobRepostings))
		}

		time.Sleep(60 * time.Minute)
	}
}

type scrapeOptions struct {
	name         string
	keywords     []string
	location     string
	workType     int64
	jobTypes     []string
	salaryRanges []string
	ageOfPosting time.Duration
}

func scrape(db *sql.DB, options scrapeOptions) (int, int, error) {
	ctx := context.Background()
	dbQueries := database.New(db)

	PROXY_HOSTNAME := os.Getenv("PROXY_HOSTNAME")
	PROXY_USERNAME := os.Getenv("PROXY_USERNAME")
	PROXY_PASSWORD := os.Getenv("PROXY_PASSWORD")

	proxyEnabled := len(PROXY_HOSTNAME) != 0

	if !proxyEnabled {
		slog.Warn("no PROXY_HOSTNAME found, proxy server disabled")
	}

	// proxy setup
	l := launcher.New()

	if proxyEnabled {
		l = l.Set(flags.ProxyServer, PROXY_HOSTNAME)
	}
	controlURL, _ := l.Launch()
	browser := rod.New()
	err := browser.ControlURL(controlURL).Connect()

	if err != nil {
		return 0, 0, fmt.Errorf("browser connection failed > %v", err)
	}

	if proxyEnabled {
		go browser.MustHandleAuth(PROXY_USERNAME, PROXY_PASSWORD)()
	}

	page, err := stealth.Page(browser)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to create stealth page > %v", err)
	}

	url, err := url.Parse("https://linkedin.com/jobs/search/")

	if err != nil {
		return 0, 0, fmt.Errorf("failed to build URL > %v", err)
	}

	q := url.Query()
	q.Set("keywords", strings.Join(options.keywords, " OR "))
	q.Set("location", options.location)
	q.Set("f_WT", fmt.Sprintf("%v", options.workType))
	q.Set("f_JT", strings.Join(options.jobTypes, ","))
	q.Set("f_SB2", strings.Join(options.salaryRanges, ","))
	q.Set("f_TPR", fmt.Sprintf("r%v", options.ageOfPosting.Seconds()))

	url.RawQuery = q.Encode()

	slog.Debug(url.String())

	err = page.Navigate(url.String())

	if err != nil {
		return 0, 0, fmt.Errorf("failed to navigate to linkedin > %v", err)
	}

	page.MustWaitStable()

	// scroll to the bottom of the page to load all virtualized resources
	scrollHeight := page.MustEval("() => document.documentElement.scrollHeight").Int()
	page.Mouse.Scroll(0.0, float64(scrollHeight), 20)

	jobPostings, err := page.Elements(".jobs-search__results-list > li")

	if err != nil {
		return 0, 0, fmt.Errorf("failed to query for job postings > %v", err)
	}

	numberOfJobPostings := len(jobPostings)
	numberOfJobRepostings := 0

	for _, jobPosting := range jobPostings {
		position, err := jobPosting.Element(".base-search-card__title")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to query for position in job posting > %v", err)
		}

		positionText, err := position.Text()

		if err != nil {
			return 0, 0, fmt.Errorf("failed to get position text from element > %v", err)
		}

		companyName, err := jobPosting.Element(".base-search-card__subtitle")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to query for company name in job posting > %v", err)
		}

		companyNameText, err := companyName.Text()

		if err != nil {
			return 0, 0, fmt.Errorf("failed to get company name text from element > %v", err)
		}

		location, err := jobPosting.Element(".job-search-card__location")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to query for location in job posting > %v", err)
		}

		locationText, err := location.Text()

		if err != nil {
			return 0, 0, fmt.Errorf("failed to get location text from element > %v", err)
		}

		postingURL, err := jobPosting.Element("a")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to query for posting url in job posting > %v", err)
		}

		postingUrlText, err := postingURL.Property("href")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to get url from element > %v", err)
		}

		relativeListingDate, err := jobPosting.Element(".job-search-card__listdate--new")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to query for relative listing date in job posting > %v", err)
		}

		relativeListingDateText, err := relativeListingDate.Text()

		if err != nil {
			return 0, 0, fmt.Errorf("failed to get url from element > %v", err)
		}

		listingDate, err := parseRelativeTime(relativeListingDateText)

		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse relative listing date > %v", err)
		}

		companyLink, err := jobPosting.Element(".base-search-card__subtitle > a")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to query for company link in job posting > %v", err)
		}

		companyLinkURL, err := companyLink.Property("href")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to get company url from element > %v", err)
		}

		parsedCompanyLinkURL, err := url.Parse(companyLinkURL.String())

		if err != nil {
			return 0, 0, fmt.Errorf("failed parsing company link url > %v", err)
		}

		segments := strings.Split(parsedCompanyLinkURL.EscapedPath(), "/")
		companySlug := segments[len(segments)-1]

		companyAvatar, err := jobPosting.Element(".base-card img")

		if err != nil {
			return 0, 0, fmt.Errorf("failed to query for company avatar in job posting > %v", err)
		}

		companyAvatarSrc, err := companyAvatar.Property("src")

		if err != nil {
			return 0, 0, fmt.Errorf("failed parsing company avatar > %v", err)
		}

		err = dbQueries.CreateCompany(ctx, database.CreateCompanyParams{ID: companySlug, Name: companyNameText, Url: companyLinkURL.String(), Avatar: sql.NullString{String: companyAvatarSrc.String(), Valid: true}})

		if err != nil {
			return 0, 0, fmt.Errorf("failed inserting company > %v", err)
		}

		err = dbQueries.CreateJobPosting(ctx, database.CreateJobPostingParams{
			Position:     positionText,
			Url:          postingUrlText.String(),
			CompanyID:    companySlug,
			LastPosted:   listingDate.UTC(),
			Location:     sql.NullString{String: locationText, Valid: true},
			LocationType: options.workType,
		})

		if err != nil {
			if sqlError, ok := err.(*sqlite.Error); ok {

				// if the posting already exists, just update the last_posted field
				if sqlError.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
					numberOfJobRepostings++

					err := dbQueries.UpdateJobPostingLastPosted(ctx, database.UpdateJobPostingLastPostedParams{
						Position:   positionText,
						CompanyID:  companySlug,
						LastPosted: listingDate.UTC(),
					})

					if err != nil {
						return 0, 0, fmt.Errorf("failed updating job posting's last_posted field > %v", err)
					}

					continue
				}
			}

			return 0, 0, fmt.Errorf("failed inserting job posting > %v", err)
		}
	}

	return numberOfJobPostings, numberOfJobRepostings, nil
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
