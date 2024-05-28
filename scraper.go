package main

import (
	"net/url"
	"strings"
)

type WorkType string

const (
	WorkTypeOnSite WorkType = "1"
	WorkTypeRemote WorkType = "2"
	WorkTypeHybrid WorkType = "3"
)

const linkedInBaseSearchURL = "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"

type ScrapeConfig struct {
	Keywords []string
	Location string
	WorkType []WorkType
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
	if len(config.WorkType) > 0 {
		q.Set("f_WT", join(config.WorkType, ","))
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
