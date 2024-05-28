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
	CompanyName string
}

type CrawledJobsPage struct {
	Jobs []CrawledJob
}

func CrawlLinkedInPage(r io.Reader) CrawledJobsPage {
	doc, _ := goquery.NewDocumentFromReader(r)

	jobElements := doc.Find("body > li")

	Jobs := make([]CrawledJob, 0, jobElements.Length())

	jobElements.Each(func(i int, s *goquery.Selection) {
		Position := strings.TrimSpace(s.Find(".base-search-card__title").Text())
		CompanyName := strings.TrimSpace(s.Find(".base-search-card__subtitle > a").Text())

		Jobs = append(Jobs, CrawledJob{
			Position,
			CompanyName,
		})
	})

	return CrawledJobsPage{
		Jobs,
	}
}
