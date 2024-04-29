package main

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/linkedincrawler"
	"github.com/m1yon/jobsummoner/internal/logger"
	_ "modernc.org/sqlite"
)

func main() {
	slog.Info("Starting...")
	logger.Init()
	godotenv.Load()

	err := linkedincrawler.Crawl()

	if err != nil {
		slog.Error("failed to crawl linkedin", tint.Err(err))
	}
}
