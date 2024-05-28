package main

import (
	"net/url"
	"strings"
)

const linkedInBaseSearchURL = "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"

type ScrapeConfig struct {
	Keywords []string
}

func BuildURL(config ScrapeConfig) string {
	url, _ := url.Parse(linkedInBaseSearchURL)

	q := url.Query()
	q.Set("keywords", strings.Join(config.Keywords, " OR "))

	url.RawQuery = q.Encode()

	return url.String()
}
