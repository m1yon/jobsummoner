package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/m1yon/jobsummoner/sql/migrations"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("libsql", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	provider, err := goose.NewProvider(database.DialectSQLite3, db, migrations.Embed)
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
	results, err := provider.Up(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("\n=== migration results  ===")
	for _, r := range results {
		log.Printf("%-3s %-2v done: %v\n", r.Source.Type, r.Source.Version, r.Duration)
	}
}
