package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/models"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestGETHomepage(t *testing.T) {

	t.Run("renders the page correctly", func(t *testing.T) {
		server := newTestServer()

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.Handler.ServeHTTP(response, request)

		assert.Equal(t, response.Code, 200)
		assertHeadingExists(t, response, "m1yon/jobsummoner")
	})

	t.Run("handles a template rendering failure", func(t *testing.T) {
		server := newTestServer()
		server.Render = func(component templ.Component, ctx context.Context, w io.Writer) error {
			return errors.New("could not render template")
		}

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.Handler.ServeHTTP(response, request)

		assert.Equal(t, 500, response.Code)
		assert.Contains(t, response.Body.String(), "Internal Server Error")
	})
}

func assertHeadingExists(t *testing.T, r *httptest.ResponseRecorder, text string) {
	t.Helper()

	doc, _ := goquery.NewDocumentFromReader(r.Body)

	found := false
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		if strings.TrimSpace(s.Text()) == text {
			found = true
		}
	})

	assert.Equal(t, true, found, fmt.Sprintf("heading with text '%v' doesn't exist", text))
}

func newTestServer() *server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db, err := sql.Open("sqlite", "./db/database.db")

	if err != nil {
		logger.Error("failed starting db")
	}

	queries := database.New(db)

	companies := &models.CompanyModel{Queries: queries}
	jobs := &models.JobModel{Queries: queries, Companies: companies}
	users := &models.UserModel{Queries: queries}

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

	jobs.CreateMany(context.Background(), jobsToCreate)

	server := newServer(logger, jobs, users, db)

	return server
}
