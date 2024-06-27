package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"path/filepath"

	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/m1yon/jobsummoner/sql/migrations"
	"github.com/pressly/goose/v3"
	gooseDB "github.com/pressly/goose/v3/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func main() {
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
		log.Fatal(err)
	}

	provider, err := goose.NewProvider(gooseDB.DialectSQLite3, db, migrations.Embed)
	if err != nil {
		log.Fatal(err)
	}
	// List migration sources the provider is aware of.
	log.Println("\n=== migration list ===")
	sources := provider.ListSources()
	for _, s := range sources {
		log.Printf("%-3s %-2v %v\n", s.Type, s.Version, filepath.Base(s.Path))
		// sql 1  00001_users_table.sql
		// go  2  00002_add_users.go
		// go  3  00003_count_users.go
	}

	// List status of migrations before applying them.
	stats, err := provider.Status(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("\n=== migration status ===")
	for _, s := range stats {
		log.Printf("%-3s %-2v %v\n", s.Source.Type, s.Source.Version, s.State)
	}

	log.Println("\n=== log migration output  ===")
	var results []*goose.MigrationResult

	if *down {
		log.Println("\n=== down migration  ===")
		result, err := provider.Down(ctx)
		results = append(results, result)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("\n=== up migrations  ===")
		results, err = provider.Up(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("\n=== migration results  ===")
	for _, r := range results {
		log.Printf("%-3s %-2v done: %v\n", r.Source.Type, r.Source.Version, r.Duration)
	}
}
