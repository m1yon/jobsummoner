package linkedin

import (
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

const (
	ErrInvalidHTML          = "HTML could not be parsed"
	ErrParsingCompanyLink   = "problem parsing company link url"
	ErrMalformedCompanyLink = "malformed company link url for parsed company link url: %v"
)

type LinkedInScraper struct {
	r      LinkedInReader
	logger *slog.Logger
}

func NewLinkedInJobScraper(r LinkedInReader, logger *slog.Logger) *LinkedInScraper {
	return &LinkedInScraper{r, logger}
}

func (l *LinkedInScraper) ScrapeJobs() ([]jobsummoner.Job, []error) {
	results := make([]jobsummoner.Job, 0)
	errs := make([]error, 0)

	for {
		reader, isLastPage, err := l.r.GetNextJobListingPage()

		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed getting page"))
			break
		}

		pageResults, err := l.scrapePage(reader)

		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed scraping page"))
			break
		}

		results = append(results, pageResults...)

		if isLastPage {
			break
		}
	}

	return results, errs
}

func (l *LinkedInScraper) scrapePage(reader io.Reader) ([]jobsummoner.Job, error) {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return []jobsummoner.Job{}, errors.Wrap(err, ErrInvalidHTML)
	}

	jobElements := doc.Find("body > li")

	Jobs := make([]jobsummoner.Job, 0, jobElements.Length())

	jobElements.Each(func(i int, s *goquery.Selection) {
		Position := strings.TrimSpace(s.Find(".base-search-card__title").Text())
		companyLinkURL, _ := s.Find(".base-search-card__subtitle > a").Attr("href")

		parsedCompanyLinkURL, err := url.Parse(companyLinkURL)

		if err != nil {
			l.logger.Error(errors.Wrap(err, ErrParsingCompanyLink).Error())
			return
		}

		segments := strings.Split(parsedCompanyLinkURL.EscapedPath(), "/")
		CompanyID := segments[len(segments)-1]

		if CompanyID == "" {
			l.logger.Error(fmt.Errorf(ErrMalformedCompanyLink, parsedCompanyLinkURL).Error())
			return
		}

		CompanyName := strings.TrimSpace(s.Find(".base-search-card__subtitle").Text())
		Location := strings.TrimSpace(s.Find(".job-search-card__location").Text())
		URL, _ := s.Find(".base-card__full-link").Attr("href")

		Jobs = append(Jobs, jobsummoner.Job{
			Position:    Position,
			CompanyID:   CompanyID,
			CompanyName: CompanyName,
			Location:    Location,
			URL:         URL,
		})
	})

	return Jobs, err
}
