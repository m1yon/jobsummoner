package http

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/stretchr/testify/assert"
)

func TestGETHomepage(t *testing.T) {
	t.Run("renders the available positions correctly", func(t *testing.T) {
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
		server.ServerHTTP(response, request)

		assert.Equal(t, response.Code, 200)

		doc, _ := goquery.NewDocumentFromReader(response.Body)
		assert.Equal(t, 1, doc.Find("h1").Length())
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
