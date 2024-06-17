package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/pkg/http"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type application struct {
	logger *slog.Logger
	db     *sql.DB
	server *http.DefaultServer
}

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	err := godotenv.Load()

	if err != nil {
		logger.Warn("no .env file found", tint.Err(err))
	}

	app := newApplication(logger)
	app.Start()
}
