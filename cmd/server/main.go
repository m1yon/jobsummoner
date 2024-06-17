package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	err := godotenv.Load()

	if err != nil {
		logger.Warn("no .env file found", tint.Err(err))
	}

	app := newApplication(logger)
	app.Start()
}
