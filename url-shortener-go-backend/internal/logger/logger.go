package logger

import (
	"log/slog"
	"os"
)

func NewLogger(env string) *slog.Logger {
	if env == "development" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func Init(env string) {
	l := NewLogger(env)
	slog.SetDefault(l)
}
