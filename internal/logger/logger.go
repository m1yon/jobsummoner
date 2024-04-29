package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func Init() {
	w := os.Stderr

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{}),
	))
}
