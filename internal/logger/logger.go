// Package logger sets up a structured JSON logger using slog.
package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	opt := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opt))
}
