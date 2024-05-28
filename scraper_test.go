package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedInPage(t *testing.T) {
	t.Run("builds correct URL given a config", func(t *testing.T) {
		got := BuildURL(ScrapeConfig{
			Keywords: []string{"react", "typescript"},
			Location: "United States",
			WorkType: []WorkType{WorkTypeRemote, WorkTypeOnSite},
		})

		assert.Equal(t, linkedInBaseSearchURL+"?f_WT=2%2C1&keywords=react+OR+typescript&location=United+States", got)
	})
	// file := os.File{}
	// got := CrawlLinkedInPage(file)

	// assert.Equal(t, LinkedInPageData{
	// 	Position:    "Software Developer",
	// 	CompanyName: "Google",
	// }, got)
}
