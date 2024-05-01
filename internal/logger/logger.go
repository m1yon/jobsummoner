package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func Init() {
	logLevel := slog.LevelInfo
	debugFlag := os.Getenv("DEBUG")

	if debugFlag == "true" {
		logLevel = slog.LevelDebug
	}

	w := os.Stderr

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level: logLevel,
		}),
	))
}
