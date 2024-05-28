package main

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type WorkType string

const (
	WorkTypeOnSite WorkType = "1"
	WorkTypeRemote WorkType = "2"
	WorkTypeHybrid WorkType = "3"
)

type JobType string

const (
	JobTypeFullTime   JobType = "F"
	JobTypePartTime   JobType = "P"
	JobTypeContract   JobType = "C"
	JobTypeTemporary  JobType = "T"
	JobTypeVolunteer  JobType = "V"
	JobTypeInternship JobType = "I"
	JobTypeOther      JobType = "O"
)

type SalaryRange string

const (
	SalaryRange40kPlus  SalaryRange = "1"
	SalaryRange60kPlus  SalaryRange = "2"
	SalaryRange80kPlus  SalaryRange = "3"
	SalaryRange100kPlus SalaryRange = "4"
	SalaryRange120kPlus SalaryRange = "5"
	SalaryRange140kPlus SalaryRange = "6"
	SalaryRange160kPlus SalaryRange = "7"
	SalaryRange180kPlus SalaryRange = "8"
	SalaryRange200kPlus SalaryRange = "9"
)

const (
	ErrParsingCompanyLink  = "problem parsing company link url: %v"
	ErrMalformedompanyLink = "malformed company link url for parsed company link url: %v"
)

const linkedInBaseSearchURL = "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"

type ScrapeConfig struct {
	Keywords    []string
	Location    string
	WorkTypes   []WorkType
	JobTypes    []JobType
	SalaryRange SalaryRange
	MaxAge      time.Duration
}

func BuildURL(config ScrapeConfig) string {
	url, _ := url.Parse(linkedInBaseSearchURL)

	q := url.Query()
	if len(config.Keywords) > 0 {
		q.Set("keywords", strings.Join(config.Keywords, " OR "))
	}
	if config.Location != "" {
		q.Set("location", config.Location)
	}
	if len(config.WorkTypes) > 0 {
		q.Set("f_WT", join(config.WorkTypes, ","))
	}
	if len(config.JobTypes) > 0 {
		q.Set("f_JT", join(config.JobTypes, ","))
	}
	if config.SalaryRange != "" {
		q.Set("f_SB2", string(config.SalaryRange))
	}
	if config.MaxAge != 0.0 {
		q.Set("f_TPR", fmt.Sprintf("r%v", config.MaxAge.Seconds()))
	}

	url.RawQuery = q.Encode()

	return url.String()
}

func join[T ~string](input []T, sep string) string {
	slice := make([]string, len(input))
	for i, v := range input {
		slice[i] = string(v)
	}

	result := strings.Join(slice, sep)

	return result
}

type CrawledJob struct {
	Position    string
	CompanyID   string
	CompanyName string
	Location    string
	URL         string
}

type CrawledJobsPage struct {
	Jobs []CrawledJob
}

func CrawlLinkedInPage(r io.Reader) (CrawledJobsPage, []error) {
	doc, _ := goquery.NewDocumentFromReader(r)

	jobElements := doc.Find("body > li")

	Jobs := make([]CrawledJob, 0, jobElements.Length())

	errors := []error{}

	jobElements.Each(func(i int, s *goquery.Selection) {
		Position := strings.TrimSpace(s.Find(".base-search-card__title").Text())
		companyLinkURL, _ := s.Find(".base-search-card__subtitle > a").Attr("href")

		parsedCompanyLinkURL, err := url.Parse(companyLinkURL)

		if err != nil {
			errors = append(errors, fmt.Errorf(ErrParsingCompanyLink, err.Error()))
			return
		}

		segments := strings.Split(parsedCompanyLinkURL.EscapedPath(), "/")
		CompanyID := segments[len(segments)-1]

		if CompanyID == "" {
			errors = append(errors, fmt.Errorf(ErrMalformedompanyLink, parsedCompanyLinkURL))
			return
		}

		CompanyName := strings.TrimSpace(s.Find(".base-search-card__subtitle").Text())
		Location := strings.TrimSpace(s.Find(".job-search-card__location").Text())
		URL, _ := s.Find(".base-card__full-link").Attr("href")

		Jobs = append(Jobs, CrawledJob{
			Position,
			CompanyID,
			CompanyName,
			Location,
			URL,
		})
	})

	return CrawledJobsPage{
		Jobs,
	}, errors
}
