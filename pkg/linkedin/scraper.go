package linkedin

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

const (
	ErrInvalidHTML         = "HTML could not be parsed"
	ErrParsingCompanyLink  = "problem parsing company link url"
	ErrMalformedompanyLink = "malformed company link url for parsed company link url: %v"
)

type LinkedInScraper struct {
	r io.Reader
}

func NewLinkedInJobScraper(r io.Reader) *LinkedInScraper {
	return &LinkedInScraper{r: r}
}

func (l *LinkedInScraper) ScrapeJobs() (jobsummoner.ScrapedJobsResults, []error) {
	errs := make([]error, 0, 1)
	doc, err := goquery.NewDocumentFromReader(l.r)

	if err != nil {
		errs = append(errs, errors.Wrap(err, ErrInvalidHTML))
		return jobsummoner.ScrapedJobsResults{}, errs
	}

	jobElements := doc.Find("body > li")

	Jobs := make([]jobsummoner.ScrapedJob, 0, jobElements.Length())

	jobElements.Each(func(i int, s *goquery.Selection) {
		Position := strings.TrimSpace(s.Find(".base-search-card__title").Text())
		companyLinkURL, _ := s.Find(".base-search-card__subtitle > a").Attr("href")

		parsedCompanyLinkURL, err := url.Parse(companyLinkURL)

		if err != nil {
			errs = append(errs, errors.Wrap(err, ErrParsingCompanyLink))
			return
		}

		segments := strings.Split(parsedCompanyLinkURL.EscapedPath(), "/")
		CompanyID := segments[len(segments)-1]

		if CompanyID == "" {
			errs = append(errs, fmt.Errorf(ErrMalformedompanyLink, parsedCompanyLinkURL))
			return
		}

		CompanyName := strings.TrimSpace(s.Find(".base-search-card__subtitle").Text())
		Location := strings.TrimSpace(s.Find(".job-search-card__location").Text())
		URL, _ := s.Find(".base-card__full-link").Attr("href")

		Jobs = append(Jobs, jobsummoner.ScrapedJob{
			Position:    Position,
			CompanyID:   CompanyID,
			CompanyName: CompanyName,
			Location:    Location,
			URL:         URL,
		})
	})

	return jobsummoner.ScrapedJobsResults{
		Jobs: Jobs,
	}, errs
}
