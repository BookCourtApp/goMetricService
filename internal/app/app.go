package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/config"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/http-server/handlers/saveHandler"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/http-server/middleware/logging"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/http-server/server"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/service"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage/clickhouse"
	RedisCache "github.com/wanna-beat-by-bit/goMetricService/internal/app/storage/redis"
)

type App struct {
	config  *config.Config
	srv     *server.Server
	storage *clickhouse.Clickhouse
	cache   *RedisCache.SessionCache
	logger  *slog.Logger
}

func New(logger *slog.Logger, conf *config.Config) (*App, error) {
	const op = "app.New"

	a := &App{}
	a.logger = logger
	a.config = conf

	appLog := a.logger.With(
		slog.String("op", op),
	)
	appLog.Debug("Building started")

	appLog.Debug("Creating database")
	db, err := clickhouse.New()
	if err != nil {
		return nil, fmt.Errorf("error while create creating database: %s", err.Error())
	}
	a.storage = db

	appLog.Debug("Initializing database")
	if err := a.storage.Init(); err != nil {
		return nil, fmt.Errorf("error while intializing clickhouse: %s", err.Error())
	}

	appLog.Debug("Creating session caching")
	a.cache, err = RedisCache.New()
	if err != nil {
		return nil, fmt.Errorf("error while initializing redis: %s", err.Error())
	}

	srvc := service.New(a.logger, a.storage, a.cache)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(logging.New(logger))
	router.Post("/test", saveHandler.New(srvc))

	a.srv = server.New(a.config, router)

	appLog.Debug("Building finished")

	return a, nil
}

func (a *App) Run() error {
	a.logger.Info("Server is up", slog.String("IP", a.srv.Network.Addr))

	if err := a.srv.Network.ListenAndServe(); err != nil {
		return fmt.Errorf("error occured while in server: %s", err.Error())
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.Stop"
	a.logger.Info("Stopping the server",
		slog.String("op", op),
	)
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := a.srv.Network.Shutdown(ctx); err != nil {
		return fmt.Errorf("error while stopping server gracefully: %s", err.Error())
	}

	return nil
}
