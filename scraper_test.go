package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedInPage(t *testing.T) {
	var tests = []struct {
		name      string
		getConfig func() ScrapeConfig
		want      string
	}{
		{
			"builds correct URL given a config",
			func() ScrapeConfig {
				config := getMockScrapeConfig()
				return config
			},
			"?f_WT=2%2C1&keywords=react+OR+typescript&location=United+States",
		},
		{
			"ignores empty Keyword field",
			func() ScrapeConfig {
				config := getMockScrapeConfig()
				config.Keywords = []string{}
				return config
			},
			"?f_WT=2%2C1&location=United+States",
		},
		{
			"ignores empty Location field",
			func() ScrapeConfig {
				config := getMockScrapeConfig()
				config.Location = ""
				return config
			},
			"?f_WT=2%2C1&keywords=react+OR+typescript",
		},
		{
			"ignores empty WorkType field",
			func() ScrapeConfig {
				config := getMockScrapeConfig()
				config.WorkType = []WorkType{}
				return config
			},
			"?keywords=react+OR+typescript&location=United+States",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildURL(tt.getConfig())
			assert.Equal(t, linkedInBaseSearchURL+tt.want, got)
		})
	}

	// file := os.File{}
	// got := CrawlLinkedInPage(file)

	// assert.Equal(t, LinkedInPageData{
	// 	Position:    "Software Developer",
	// 	CompanyName: "Google",
	// }, got)
}

func getMockScrapeConfig() ScrapeConfig {
	config := ScrapeConfig{
		Keywords: []string{"react", "typescript"},
		Location: "United States",
		WorkType: []WorkType{WorkTypeRemote, WorkTypeOnSite},
	}

	return config
}
