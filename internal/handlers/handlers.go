package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Templates struct {
	templates *template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/**/*.html")),
	}
}

type handlersConfig struct {
	DB *database.Queries
}

func NewHandlerMux(db *sql.DB) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	dbQueries := database.New(db)

	cfg := handlersConfig{
		DB: dbQueries,
	}

	mux.HandleFunc("/", cfg.rootHandler)

	mux.HandleFunc("GET /user_job_postings", cfg.getUserJobPostingsHandler)
	mux.HandleFunc("PATCH /user_job_postings/{jobPostingID}", cfg.patchUserJobPostingsHandler)

	mux.Handle("/metrics/", promhttp.Handler())

	return mux, nil
}
