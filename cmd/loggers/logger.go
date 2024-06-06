package loggers

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

// SetupLogger создание нового логгера.
func SetupLogger(env string) *Logger {
	var logger *slog.Logger
	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return &Logger{Logger: logger}
}
