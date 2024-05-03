package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/internal/handlers"
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

	dbQueries := database.New(db)

	// create our seed user
	_, err = dbQueries.GetUser(ctx, 1)
	if err != nil {
		dbQueries.CreateUser(ctx)
	}

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
