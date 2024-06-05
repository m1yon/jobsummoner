package main

import (
	"log/slog"
	"os"

	"github.com/m1yon/jobsummoner/pkg/http"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	server := http.NewDefaultServer(logger)
	server.ListenAndServe(":3000")
}
