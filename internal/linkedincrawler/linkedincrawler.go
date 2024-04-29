package linkedincrawler

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/stealth"
)

func Crawl() error {

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
		return fmt.Errorf("browser connection failed > %v", err)
	}

	if proxyEnabled {
		go browser.MustHandleAuth(PROXY_USERNAME, PROXY_PASSWORD)()
	}

	page, err := stealth.Page(browser)

	if err != nil {
		return fmt.Errorf("failed to create stealth page > %v", err)
	}

	keywords := []string{"typescript", "react"}
	location := "United States"
	workTypes := []string{"2"}    // 2 = remote
	jobTypes := []string{"F"}     // F = fulltime
	salaryRanges := []string{"5"} // 5 = $120,000+
	ageOfPosting := 24 * time.Hour
	url, err := url.Parse("https://linkedin.com/jobs/search/")

	if err != nil {
		return fmt.Errorf("failed to build URL > %v", err)
	}

	q := url.Query()
	q.Set("keywords", strings.Join(keywords, " OR "))
	q.Set("location", location)
	q.Set("f_WT", strings.Join(workTypes, ","))
	q.Set("f_JT", strings.Join(jobTypes, ","))
	q.Set("f_SB2", strings.Join(salaryRanges, ","))
	q.Set("f_TPR", fmt.Sprintf("r%v", ageOfPosting.Seconds()))

	url.RawQuery = q.Encode()

	err = page.Navigate(url.String())

	if err != nil {
		return fmt.Errorf("failed to navigate to linkedin > %v", err)
	}

	page.MustWaitStable()

	jobPostings, err := page.Elements(".jobs-search__results-list > li")

	if err != nil {
		return fmt.Errorf("failed to query for job postings > %v", err)
	}

	for _, jobPosting := range jobPostings {
		position, err := jobPosting.Element(".base-search-card__title")

		if err != nil {
			return fmt.Errorf("failed to query for position in job posting > %v", err)
		}

		positionText, err := position.Text()

		if err != nil {
			return fmt.Errorf("failed to get position text from element > %v", err)
		}

		companyName, err := jobPosting.Element(".base-search-card__subtitle")

		if err != nil {
			return fmt.Errorf("failed to query for company name in job posting > %v", err)
		}

		companyNameText, err := companyName.Text()

		if err != nil {
			return fmt.Errorf("failed to get company name text from element > %v", err)
		}

		postingURL, err := jobPosting.Element("a")

		if err != nil {
			return fmt.Errorf("failed to query for posting url in job posting > %v", err)
		}

		urlText, err := postingURL.Property("href")

		if err != nil {
			return fmt.Errorf("failed to get url from element > %v", err)
		}

		companyLink, err := jobPosting.Element(".base-search-card__subtitle > a")

		if err != nil {
			return fmt.Errorf("failed to query for company link in job posting > %v", err)
		}

		companyLinkURL, err := companyLink.Property("href")

		if err != nil {
			return fmt.Errorf("failed to get company url from element > %v", err)
		}

		parsedCompanyLinkURL, err := url.Parse(companyLinkURL.String())

		if err != nil {
			return fmt.Errorf("failed parsing company link url > %v", err)
		}

		segments := strings.Split(parsedCompanyLinkURL.EscapedPath(), "/")
		companySlug := segments[len(segments)-1]

		fmt.Println(positionText, companyNameText, urlText, companyLinkURL, companySlug)
	}

	return nil
}
