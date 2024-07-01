package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1yon/jobsummoner/internal/models"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestHomepage(t *testing.T) {
	jobsToCreate := []models.Job{
		{
			Position:      "Software Developer",
			URL:           "https://linkedin.com/jobs/1",
			Location:      "San Francisco",
			SourceID:      "linkedin",
			CompanyID:     "/google",
			CompanyName:   "Google",
			CompanyAvatar: "https://upload.wikimedia.org/wikipedia/commons/c/c1/Google_%22G%22_logo.svg",
		},
		{
			Position:      "Manager",
			URL:           "https://linkedin.com/jobs/2",
			Location:      "Seattle",
			SourceID:      "linkedin",
			CompanyID:     "/microsoft",
			CompanyName:   "Microsoft",
			CompanyAvatar: "https://blogs.microsoft.com/wp-content/uploads/prod/2012/08/8867.Microsoft_5F00_Logo_2D00_for_2D00_screen.jpg",
		},
	}

	t.Run("renders the page correctly", func(t *testing.T) {
		app := newTestApplication(t)

		ts := newTestServer(t, app.routes())
		defer ts.Close()

		app.jobs.CreateMany(context.Background(), jobsToCreate)

		code, _, body := ts.get(t, "/")

		assert.Equal(t, 200, code)
		assertHeadingExists(t, body, "m1yon/jobsummoner")
	})
}

func assertHeadingExists(t *testing.T, body string, text string) {
	t.Helper()

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))

	found := false
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		if strings.TrimSpace(s.Text()) == text {
			found = true
		}
	})

	assert.Equal(t, true, found, fmt.Sprintf("heading with text '%v' doesn't exist", text))
}
