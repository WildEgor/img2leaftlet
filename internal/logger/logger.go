// Package logger ...
package logger

import (
	"log/slog"
	"os"
)

// Init logger
func Init() {
	var opts = &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	slog.SetDefault(slog.New(handler))
}
