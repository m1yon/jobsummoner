package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/m1yon/jobsummoner/pkg/http"
)

type application struct {
	logger *slog.Logger
	db     *sql.DB
	server *http.Server
	config *config
}

func newApplication(logger *slog.Logger) *application {
	config := getConfigFromFlags()

	db, err := openDB(logger, config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	server := newServer(logger, db)

	return &application{logger: logger, db: db, server: server, config: config}
}

func (a *application) Start() {
	a.server.Start(":3000")
}
