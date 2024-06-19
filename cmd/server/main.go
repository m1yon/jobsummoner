package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	app := newApplication(logger)
	app.Start()
}
