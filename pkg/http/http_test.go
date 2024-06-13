package http

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
	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestGETHomepage(t *testing.T) {
	jobsToCreate := []jobsummoner.Job{
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
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		db, _ := sqlitedb.NewInMemoryDB()

		companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
		companyService := company.NewDefaultCompanyService(companyRepository)
		jobRepository := sqlitedb.NewSqliteJobRepository(db)
		jobService := job.NewDefaultJobService(jobRepository, companyService)

		jobService.CreateJobs(ctx, jobsToCreate)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server := NewDefaultServer(logger, jobService)
		server.ServerHTTP(response, request)

		assert.Equal(t, response.Code, 200)

		doc, _ := goquery.NewDocumentFromReader(response.Body)

		assertHeadingExists(t, doc, "m1yon/jobsummoner")
	})

	t.Run("handles a template rendering failure", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		db, err := sql.Open("sqlite", "./db/database.db")

		if err != nil {
			logger.Error("failed starting db")
		}

		companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
		companyService := company.NewDefaultCompanyService(companyRepository)
		jobRepository := sqlitedb.NewSqliteJobRepository(db)
		jobService := job.NewDefaultJobService(jobRepository, companyService)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server := NewDefaultServer(logger, jobService)
		server.Render = func(component templ.Component, ctx context.Context, w io.Writer) error {
			return errors.New("could not render template")
		}
		server.ServerHTTP(response, request)

		assert.Equal(t, 500, response.Code)
		assert.Equal(t, "", response.Body.String())
	})
}

func assertHeadingExists(t *testing.T, doc *goquery.Document, text string) {
	t.Helper()

	found := false
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(i int, s *goquery.Selection) {
		if strings.TrimSpace(s.Text()) == text {
			found = true
		}
	})

	assert.Equal(t, true, found, fmt.Sprintf("heading with text '%v' doesn't exist", text))
}
