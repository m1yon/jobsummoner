package linkedin

import (
	"fmt"
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestLinkedInScraper(t *testing.T) {
	t.Run("scrapes the jobs correctly", func(t *testing.T) {
		mockReader := NewMockLinkedInReader("./test-helpers/li-job-listings.html")
		scraper := NewLinkedInJobScraper(mockReader)
		got, errs := scraper.ScrapeJobs()

		assert.Equal(t, 0, len(errs))
		assert.Equal(t, jobsummoner.ScrapedJobsResults{
			Jobs: []jobsummoner.Job{
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
		mockReader := NewMockLinkedInReader("./test-helpers/li-job-listings_bad-company-url.html")
		scraper := NewLinkedInJobScraper(mockReader)
		got, errs := scraper.ScrapeJobs()

		assert.Equal(t, errs, []error{
			fmt.Errorf(ErrMalformedompanyLink, "fda&=+!-//"),
		})
		assert.Equal(t, jobsummoner.ScrapedJobsResults{
			Jobs: []jobsummoner.Job{
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
