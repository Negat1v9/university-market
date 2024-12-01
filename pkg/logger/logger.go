package logger

import (
	"log/slog"
	"os"
)

const (
	envProd  = "prod"
	envLocal = "local"
)

func New(env string) *slog.Logger {
	var log *slog.Logger

	switch env {

	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
		log.Info("logger", slog.Any("level", env))

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
		)
		log.Info("logger", slog.Any("level", env))
	default:
		panic("environment variable for logger not specified")
	}
	return log
}
