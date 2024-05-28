package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedInURLBuilder(t *testing.T) {
	var tests = []struct {
		name      string
		getConfig func() ScrapeConfig
		want      string
	}{
		{
			"Keywords field",
			func() ScrapeConfig {
				config := ScrapeConfig{Keywords: []string{"react", "typescript"}}
				return config
			},
			"?keywords=react+OR+typescript",
		},
		{
			"Location field",
			func() ScrapeConfig {
				config := ScrapeConfig{Location: "United States"}
				return config
			},
			"?location=United+States",
		},
		{
			"WorkTypes field",
			func() ScrapeConfig {
				config := ScrapeConfig{WorkTypes: []WorkType{WorkTypeRemote, WorkTypeOnSite}}
				return config
			},
			"?f_WT=2%2C1",
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
