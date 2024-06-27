package database

import (
	"database/sql"
	"os"

	"github.com/pkg/errors"
)

const (
	ErrOpeningDB = "problem opening db"
	ErrPingingDB = "db did not respond to ping"
	ErrDSNNotSet = "dsn not provided"
)

func NewFileDB(opener ConnectionOpener) (*sql.DB, error) {
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

	return opener.Open("sqlite", localDbPath)
}

func NewTursoDB(dsn string, opener ConnectionOpener) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New(ErrDSNNotSet)
	}

	db, err := opener.Open("libsql", dsn)

	if err != nil {
		return nil, errors.Wrap(err, ErrOpeningDB)
	}

	return db, nil
}
