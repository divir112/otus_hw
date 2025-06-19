package logger

import (
	"io"
	"log/slog"
)

func New(level string, file io.Writer) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case "INFO":
		slogLevel = slog.LevelInfo
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "ERROR":
		slogLevel = slog.LevelError
	case "WARN":
		slogLevel = slog.LevelWarn
	}
	return slog.New(slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slogLevel,
	}))
}
