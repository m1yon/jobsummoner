package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedInPage(t *testing.T) {
	t.Run("builds correct URL given a config", func(t *testing.T) {
		got := BuildURL(ScrapeConfig{
			Keywords: []string{"react", "typescript"},
		})

		assert.Equal(t, linkedInBaseSearchURL+"?keywords=react+OR+typescript", got)
	})
	// file := os.File{}
	// got := CrawlLinkedInPage(file)

	// assert.Equal(t, LinkedInPageData{
	// 	Position:    "Software Developer",
	// 	CompanyName: "Google",
	// }, got)
}
