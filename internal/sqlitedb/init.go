package sqlitedb

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/pkg/errors"
)

const (
	ErrOpeningDB         = "problem opening db"
	ErrPingingDB         = "db did not respond to ping"
	ErrDatabaseURLNotSet = "DATABASE_URL not set"
)

func NewDB(logger *slog.Logger, open func(driverName string, dataSourceName string) (*sql.DB, error)) (*sql.DB, error) {
	useLocalDB := os.Getenv("LOCAL_DB")

	if useLocalDB == "true" {
		logger.Info("using local database")
		return NewFileDB()
	}

	return NewTursoDB(open)
}

func NewTursoDB(open func(driverName string, dataSourceName string) (*sql.DB, error)) (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		return nil, errors.New(ErrDatabaseURLNotSet)
	}

	db, err := open("libsql", dbURL)

	if err != nil {
		return nil, errors.Wrap(err, ErrOpeningDB)
	}

	return db, nil
}
