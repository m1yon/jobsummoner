package handlers

import (
	"database/sql"
	"net/http"

	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
	mux.HandleFunc("PATCH /user_job_postings/{jobPostingID}", cfg.putUserJobPostingsHandler)

	mux.Handle("/metrics/", promhttp.Handler())

	return mux, nil
}
