package sqlitedb

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
)

func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}

func NewFileDB() (*sql.DB, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		return nil, errors.Wrap(err, "error getting working directory")
	}

	migrationsDir := workingDir + "/sql/migrations"
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

	return migrateLocalDB(localDbPath, migrationsDir)
}

func NewInMemoryDB() (*sql.DB, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		return nil, errors.Wrap(err, "error getting working directory")
	}

	migrationsDir := workingDir + "/sql/migrations"
	fmt.Println("migrationsDir", migrationsDir)

	return migrateLocalDB(":memory:", migrationsDir)
}

func migrateLocalDB(dataSourceName string, migrationsDir string) (*sql.DB, error) {
	ctx := context.Background()
	db, err := goose.OpenDBWithDriver("sqlite3", dataSourceName)

	if err != nil {
		return nil, errors.Wrap(err, "failed to open DB")
	}

	if err := goose.RunContext(ctx, "up", db, migrationsDir); err != nil {
		return nil, errors.Wrap(err, "goose error")
	}

	return db, nil
}
