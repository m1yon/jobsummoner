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
		return NewFileDB(open)
	}

	return NewTursoDB(open)
}

func NewFileDB(open func(driverName string, dataSourceName string) (*sql.DB, error)) (*sql.DB, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		return nil, errors.Wrap(err, "error getting working directory")
	}

	// migrationsDir := workingDir + "/sql/migrations"
	localDbDir := workingDir + "/db"
	localDbPath := localDbDir + "/database.db"

	if _, err := os.Stat(localDbDir); os.IsNotExist(err) {
		err := os.Mkdir(localDbDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	_, err = os.OpenFile(localDbPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return open("sqlite", localDbPath)
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
