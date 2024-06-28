package database

import (
	"context"
	"database/sql"
	"os"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
)

func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}

func NewInMemoryDB() (*sql.DB, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		return nil, errors.Wrap(err, "error getting working directory")
	}

	migrationsDir := workingDir + "/sql/migrations"

	return migrateLocalDB(":memory:", migrationsDir)
}

func migrateLocalDB(dataSourceName string, migrationsDir string) (*sql.DB, error) {
	ctx := context.Background()
	goose.SetLogger(goose.NopLogger())
	db, err := goose.OpenDBWithDriver("sqlite3", dataSourceName)

	if err != nil {
		return nil, errors.Wrap(err, "failed to open DB")
	}

	if err := goose.RunContext(ctx, "up", db, migrationsDir, "> /dev/null"); err != nil {
		return nil, errors.Wrap(err, "goose error")
	}

	return db, nil
}
