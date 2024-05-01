package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/linkedincrawler"
	"github.com/m1yon/jobsummoner/internal/logger"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	logger.Init()
	godotenv.Load()
	port := "3000"

	db, err := sql.Open("sqlite", "./db/database.db")

	if err != nil {
		slog.Error("problem initializing handler mux", err)
		return
	}
	go linkedincrawler.ScrapeLoop(db)

	mux := http.NewServeMux()
	dbQueries := database.New(db)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("cmd/jobsummoner/index.html")

		if err != nil {
			slog.Error("could not parse template", tint.Err(err))
		}

		JobPostings, err := dbQueries.GetJobPostings(ctx)

		if err != nil {
			slog.Error("failed to query job postings", tint.Err(err))
		}

		type FormattedJobPosting struct {
			database.GetJobPostingsRow
			TimeAgo string
		}

		formattedJobPostings := make([]FormattedJobPosting, 0, len(JobPostings))

		for _, jobPosting := range JobPostings {
			formattedJobPostings = append(formattedJobPostings, FormattedJobPosting{GetJobPostingsRow: jobPosting, TimeAgo: timeAgo(jobPosting.LastPosted)})
		}

		homepage := struct {
			JobPostings []FormattedJobPosting
		}{
			formattedJobPostings,
		}

		err = t.Execute(w, homepage)

		if err != nil {
			slog.Error("could not execute template", tint.Err(err))
		}
	})

	corsMux := middlewareCors(mux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	slog.Info("Running server", slog.String("port", port))
	server.ListenAndServe()
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func timeAgo(from time.Time) string {
	now := time.Now()
	diff := now.Sub(from)

	if diff < time.Minute {
		if int(diff.Seconds()) == 1 {
			return fmt.Sprintf("%d second ago", int(diff.Seconds()))
		}
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	} else if diff < time.Hour {
		if int(diff.Minutes()) == 1 {
			return fmt.Sprintf("%d minute ago", int(diff.Minutes()))
		}
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	} else if diff < time.Hour*24 {
		if int(diff.Hours()) == 1 {
			return fmt.Sprintf("%d hour ago", int(diff.Hours()))
		}
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	} else if diff < time.Hour*24*30 {
		days := diff / (time.Hour * 24)
		if days == 1 {
			return fmt.Sprintf("%d day ago", int(days))
		}
		return fmt.Sprintf("%d days ago", int(days))
	} else if diff < time.Hour*24*365 {
		months := diff / (time.Hour * 24 * 30)
		if months == 1 {
			return fmt.Sprintf("%d month ago", int(months))
		}
		return fmt.Sprintf("%d months ago", int(months))
	}
	years := diff / (time.Hour * 24 * 365)
	if years == 1 {
		return fmt.Sprintf("%d year ago", int(years))
	}
	return fmt.Sprintf("%d years ago", int(years))
}
