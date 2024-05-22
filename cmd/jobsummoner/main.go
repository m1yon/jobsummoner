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

	corsMux := handlers.MiddlewareCors(mux)
	loggingMux := handlers.MiddlewareLogging(corsMux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: loggingMux,
	}

	slog.Info("Running server", slog.String("port", port))
	server.ListenAndServe()
}

func seedDB(db *sql.DB) {
	ctx := context.Background()

	dbQueries := database.New(db)

	// create our seed user
	_, err := dbQueries.GetUser(ctx, 1)
	if err != nil {
		dbQueries.CreateUser(ctx)
	}

	scrapes, err := dbQueries.GetAllScrapes(ctx)
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

		keywords := []string{"typescript", "react"}
		for _, keyword := range keywords {
			err = dbQueries.AddKeywordToScrape(ctx, database.AddKeywordToScrapeParams{
				ScrapeID: 1,
				Keyword:  keyword,
			})

			if err != nil {
				slog.Error("failed inserting seed scrape keyword", tint.Err(err))
				return
			}
		}

		positionBlacklistedWords := []string{"manager", "executive", "staff", "principal", "architect", "react native", "lead", "designer", "python", "jr", "technologist", "director", "clearance"}
		for _, positionBlacklistedWord := range positionBlacklistedWords {
			err = dbQueries.AddPositionBlacklistedWordToScrape(ctx, database.AddPositionBlacklistedWordToScrapeParams{
				ScrapeID:        1,
				BlacklistedWord: positionBlacklistedWord,
			})

			if err != nil {
				slog.Error("failed inserting seed scrape position blacklisted word", tint.Err(err))
				return
			}
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

		for _, keyword := range keywords {
			err = dbQueries.AddKeywordToScrape(ctx, database.AddKeywordToScrapeParams{
				ScrapeID: 2,
				Keyword:  keyword,
			})

			if err != nil {
				slog.Error("failed inserting seed scrape keyword", tint.Err(err))
				return
			}
		}

		for _, positionBlacklistedWord := range positionBlacklistedWords {
			err = dbQueries.AddPositionBlacklistedWordToScrape(ctx, database.AddPositionBlacklistedWordToScrapeParams{
				ScrapeID:        2,
				BlacklistedWord: positionBlacklistedWord,
			})

			if err != nil {
				slog.Error("failed inserting seed scrape position blacklisted word", tint.Err(err))
				return
			}
		}
	}
}
