package sqlitedb

import (
	"context"
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
)

const (
	dbstring = ":memory:"
	dir      = "../../sql/migrations"
)

func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}

func NewTestDB() *sql.DB {
	ctx := context.Background()
	db, err := goose.OpenDBWithDriver("sqlite3", dbstring)

	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
		return &sql.DB{}
	}

	if err := goose.RunContext(ctx, "up", db, dir); err != nil {
		log.Fatalf("goose erro: %v", err)
		return &sql.DB{}
	}

	return db
}
