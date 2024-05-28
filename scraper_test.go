package jobsummoner

import (
	"fmt"
	"os"
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

}

func TestLinkedInScraper(t *testing.T) {
	t.Run("scrapes the jobs correctly", func(t *testing.T) {
		file, err := os.Open("./tests/li-job-listings.html")

		if err != nil {
			t.Fatal("could not open file")
		}

		got, errs := CrawlLinkedInPage(file)

		assert.Equal(t, 0, len(errs))
		assert.Equal(t, CrawledJobsPage{
			Jobs: []CrawledJob{
				{
					Position:    "Full Stack Engineer",
					CompanyID:   "venchrpartners",
					CompanyName: "Venchr",
					Location:    "San Francisco Bay Area",
					URL:         "https://www.linkedin.com/jobs/view/full-stack-engineer-at-venchr-3936836372?position=1&pageNum=0&refId=fsDMYm%2BoJB2zdtWm%2FhnZ3g%3D%3D&trackingId=lBDr2giv8wtQWgYNy9o7zA%3D%3D&trk=public_jobs_jserp-result_search-card",
				},
				{
					Position:    "Software Engineer II (Frontend) - Seller Experience",
					CompanyID:   "stubhub",
					CompanyName: "StubHub",
					Location:    "Los Angeles, CA",
					URL:         "https://www.linkedin.com/jobs/view/software-engineer-ii-frontend-seller-experience-at-stubhub-3916280897?position=2&pageNum=0&refId=fsDMYm%2BoJB2zdtWm%2FhnZ3g%3D%3D&trackingId=dbosG%2Ftu2ZxD8zDrXnrTWw%3D%3D&trk=public_jobs_jserp-result_search-card",
				},
				{
					Position:    "Senior Frontend Developer",
					CompanyID:   "trilogy-international-ltd",
					CompanyName: "Trilogy International",
					Location:    "South San Francisco, CA",
					URL:         "https://www.linkedin.com/jobs/view/senior-frontend-developer-at-trilogy-international-3936896077?position=3&pageNum=0&refId=fsDMYm%2BoJB2zdtWm%2FhnZ3g%3D%3D&trackingId=YzFScrnDUB9kJlZ%2FHtqdiQ%3D%3D&trk=public_jobs_jserp-result_search-card",
				},
			},
		}, got)
	})

	t.Run("handles invalid company URLs", func(t *testing.T) {
		file, err := os.Open("./tests/li-job-listings_bad-company-url.html")

		if err != nil {
			t.Fatal("could not open file")
		}

		got, errs := CrawlLinkedInPage(file)

		assert.Equal(t, errs, []error{
			fmt.Errorf(ErrMalformedompanyLink, "fda&=+!-//"),
		})
		assert.Equal(t, CrawledJobsPage{
			Jobs: []CrawledJob{
				{
					Position:    "Software Engineer II (Frontend) - Seller Experience",
					CompanyID:   "stubhub",
					CompanyName: "StubHub",
					Location:    "Los Angeles, CA",
					URL:         "https://www.linkedin.com/jobs/view/software-engineer-ii-frontend-seller-experience-at-stubhub-3916280897?position=2&pageNum=0&refId=fsDMYm%2BoJB2zdtWm%2FhnZ3g%3D%3D&trackingId=dbosG%2Ftu2ZxD8zDrXnrTWw%3D%3D&trk=public_jobs_jserp-result_search-card",
				},
				{
					Position:    "Senior Frontend Developer",
					CompanyID:   "trilogy-international-ltd",
					CompanyName: "Trilogy International",
					Location:    "South San Francisco, CA",
					URL:         "https://www.linkedin.com/jobs/view/senior-frontend-developer-at-trilogy-international-3936896077?position=3&pageNum=0&refId=fsDMYm%2BoJB2zdtWm%2FhnZ3g%3D%3D&trackingId=YzFScrnDUB9kJlZ%2FHtqdiQ%3D%3D&trk=public_jobs_jserp-result_search-card",
				},
			},
		}, got)
	})
}
