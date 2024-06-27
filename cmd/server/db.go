package main

import (
	"database/sql"
	"log/slog"

	"github.com/m1yon/jobsummoner/internal/database"
)

func openDB(logger *slog.Logger, config *config) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if config.useLocalDB {
		db, err = database.NewFileDB(&database.SqlConnectionOpener{})
		logger.Info("using local DB")
	} else {
		db, err = database.NewTursoDB(config.dsn, &database.SqlConnectionOpener{})
		logger.Info("using remote DB")
	}

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
