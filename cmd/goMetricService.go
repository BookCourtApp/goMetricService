package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/wanna-beat-by-bit/goMetricService/internal/app"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	conf := config.MustLoad()
	logger := setupLogger(conf.Env)

	logger.Info(
		"Logger enabled",
		slog.String("env", conf.Env),
	)
	logger.Debug("debug messages is enabled")

	a, err := app.New(logger, conf)
	if err != nil {
		logger.Error("fatal error occured", slog.String("error", err.Error()))
		panic("can't build application")
	}

	logger.Info("starting application")
	if err := a.Run(); err != nil {
		logger.Error("fatal error occured",
			slog.String("error", err.Error()),
		)
	}

	if err := a.Stop(context.Background()); err != nil {
		logger.Error("error while stopping server",
			slog.String("error", err.Error()),
		)
	}

	logger.Info("application exit")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
