package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/handlers"
	"github.com/m1yon/jobsummoner/internal/linkedincrawler"
	"github.com/m1yon/jobsummoner/internal/logger"
	_ "modernc.org/sqlite"
)

func main() {

	logger.Init()
	godotenv.Load()
	port := "3000"

	db, err := sql.Open("sqlite", "./db/database.db")

	if err != nil {
		slog.Error("problem initializing handler mux", err)
		return
	}

	seedDB(db)

	go linkedincrawler.ScrapeLoop(db)

	mux, err := handlers.NewHandlerMux(db)

	if err != nil {
		slog.Error("problem initializing handler mux", err)
		return
	}

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

func seedDB(db *sql.DB) {
	ctx := context.Background()

	dbQueries := database.New(db)

	// create our seed user
	_, err := dbQueries.GetUser(ctx, 1)
	if err != nil {
		dbQueries.CreateUser(ctx)
	}

	scrapes, err := dbQueries.GetAllScrapesWithKeywords(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "converting NULL to int64 is unsupported") {
			slog.Error("failed querying for all scrapes", tint.Err(err))
			return
		}
	}

	if len(scrapes) == 0 {
		err := dbQueries.CreateScrape(ctx, database.CreateScrapeParams{
			Name:     "Remote Roles",
			Location: "United States",
			WorkType: 2,
			UserID:   1,
		})

		if err != nil {
			slog.Error("failed inserting seed scrape", tint.Err(err))
			return
		}

		err = dbQueries.AddKeywordToScrape(ctx, database.AddKeywordToScrapeParams{
			ScrapeID: 1,
			Keyword:  "typescript",
		})

		if err != nil {
			slog.Error("failed inserting seed scrape keyword", tint.Err(err))
			return
		}

		err = dbQueries.AddKeywordToScrape(ctx, database.AddKeywordToScrapeParams{
			ScrapeID: 1,
			Keyword:  "react",
		})

		if err != nil {
			slog.Error("failed inserting seed scrape keyword", tint.Err(err))
			return
		}

		err = dbQueries.CreateScrape(ctx, database.CreateScrapeParams{
			Name:     "Colorado Hybrid Roles",
			Location: "Colorado, United States",
			WorkType: 3,
			UserID:   1,
		})

		if err != nil {
			slog.Error("failed inserting seed scrape", tint.Err(err))
			return
		}

		err = dbQueries.AddKeywordToScrape(ctx, database.AddKeywordToScrapeParams{
			ScrapeID: 2,
			Keyword:  "typescript",
		})

		if err != nil {
			slog.Error("failed inserting seed scrape keyword", tint.Err(err))
			return
		}

		err = dbQueries.AddKeywordToScrape(ctx, database.AddKeywordToScrapeParams{
			ScrapeID: 2,
			Keyword:  "react",
		})

		if err != nil {
			slog.Error("failed inserting seed scrape keyword", tint.Err(err))
			return
		}

	}
}
