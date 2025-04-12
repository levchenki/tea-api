package slogx

import (
	"github.com/levchenki/tea-api/internal/config"
	"log/slog"
	"os"
)

func Setup(env config.Environment) *slog.Logger {
	var logger *slog.Logger
	switch env {

	case config.EnvLocal:
		logger = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		))
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		))
	}

	return logger
}
