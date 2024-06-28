package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/sql/migrations"
	"github.com/pressly/goose/v3"
	gooseDB "github.com/pressly/goose/v3/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{}))

	ctx := context.Background()

	useLocalDB := flag.Bool("local-db", true, "Use a local sqlite DB")
	dsn := flag.String("dsn", "", "Database connection string")
	down := flag.Bool("down", false, "Runs a down migration")

	flag.Parse()

	var db *sql.DB
	var err error

	if *useLocalDB {
		db, err = database.NewFileDB(&database.SqlConnectionOpener{})
	} else {
		db, err = database.NewTursoDB(*dsn, &database.SqlConnectionOpener{})
	}

	if err != nil {
		logger.Error("failed to init db", tint.Err(err))
		os.Exit(1)
	}

	provider, err := goose.NewProvider(gooseDB.DialectSQLite3, db, migrations.Embed)
	if err != nil {
		logger.Error("failed to init provider", tint.Err(err))
		os.Exit(1)
	}

	// List migration sources the provider is aware of.
	sources := provider.ListSources()
	for _, s := range sources {
		logger.Info("migration source found", "type", s.Type, "version", s.Version, "path", filepath.Base(s.Path))
	}

	// List status of migrations before applying them.
	stats, err := provider.Status(ctx)
	if err != nil {
		logger.Error("failed get migrations status", tint.Err(err))
		os.Exit(1)
	}
	for _, s := range stats {
		logger.Info("current migration status", "type", s.Source.Type, "version", s.Source.Version, "state", s.State)
	}

	var results []*goose.MigrationResult
	if *down {
		result, err := provider.Down(ctx)
		results = append(results, result)

		if err != nil {
			logger.Error("failed down migrations", tint.Err(err))
			os.Exit(1)
		}
	} else {
		results, err = provider.Up(ctx)
		if err != nil {
			logger.Error("failed up migrations", tint.Err(err))
			os.Exit(1)
		}
	}

	for _, r := range results {
		logger.Info("migration successful", "type", r.Source.Type, "version", r.Source.Version, "duration", r.Duration)
	}
}
