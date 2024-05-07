package handlers

import (
	"database/sql"
	"io"
	"net/http"
	"text/template"

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
	DB       *database.Queries
	Renderer *Templates
}

func (t *Templates) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewHandlerMux(db *sql.DB) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	dbQueries := database.New(db)

	cfg := handlersConfig{
		DB:       dbQueries,
		Renderer: NewTemplates(),
	}

	mux.HandleFunc("/", cfg.rootHandler)
	mux.HandleFunc("PATCH /user_job_postings/{jobPostingID}", cfg.putUserJobPostingsHandler)

	mux.Handle("/metrics/", promhttp.Handler())

	return mux, nil
}
