package main

import (
	"database/sql"
	"log/slog"

	"github.com/m1yon/jobsummoner/internal/sqlitedb"
)

func openDB(logger *slog.Logger) (*sql.DB, error) {
	db, err := sqlitedb.NewDB(logger, &sqlitedb.SqlConnectionOpener{})

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
