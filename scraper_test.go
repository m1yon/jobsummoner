package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLinkedInURLBuilder(t *testing.T) {
	tests := []struct {
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
		{
			"JobTypes field",
			func() ScrapeConfig {
				config := ScrapeConfig{JobTypes: []JobType{JobTypeFullTime, JobTypeOther}}
				return config
			},
			"?f_JT=F%2CO",
		},
		{
			"SalaryRange field",
			func() ScrapeConfig {
				config := ScrapeConfig{SalaryRange: SalaryRange160kPlus}
				return config
			},
			"?f_SB2=7",
		},
		{
			"MaxAge field",
			func() ScrapeConfig {
				config := ScrapeConfig{MaxAge: time.Hour * 24}
				return config
			},
			"?f_TPR=r86400",
		},
		{
			"All fields",
			func() ScrapeConfig {
				config := ScrapeConfig{
					Keywords:    []string{"go", "templ"},
					Location:    "Africa",
					WorkTypes:   []WorkType{WorkTypeHybrid},
					JobTypes:    []JobType{JobTypeFullTime, JobTypeOther},
					SalaryRange: SalaryRange200kPlus,
					MaxAge:      time.Hour * 12,
				}
				return config
			},
			"?f_JT=F%2CO&f_SB2=9&f_TPR=r43200&f_WT=3&keywords=go+OR+templ&location=Africa",
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
