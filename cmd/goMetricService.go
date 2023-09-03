package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"time"

	//"github.com/wanna-beat-by-bit/goMetricService/internal/app"
	//"github.com/wanna-beat-by-bit/goMetricService/internal/app/config"

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
	logger.Debug("debug messages enabled")

	a, err := New(logger, conf)
	if err != nil {
		logger.Error("fatal error occured", slog.String("error", err.Error()))
		panic("can't build application")
	}

	sysExit := make(chan os.Signal, 1)
	signal.Notify(sysExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if err := a.Run(); err != nil {
		logger.Error("fatal error occured",
			slog.String("error", err.Error()),
		)
	}
	logger.Info("application started")

	<-sysExit

	logger.Info("Got syscall to exit")
	if err := a.Stop(context.Background()); err != nil {
		logger.Error("error while exiting",
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

type App struct {
	config  *config.Config
	srv     *server.Server
	storage *clickhouse.Clickhouse
	cache   *RedisCache.SessionCache
	logger  *slog.Logger
}

func New(logger *slog.Logger, conf *config.Config) (*App, error) {
	const op = "app.New"
	var err error

	a := &App{}
	a.logger = logger
	a.config = conf

	appLog := a.logger.With(
		slog.String("op", op),
	)
	appLog.Debug("Building started")

	appLog.Debug("Creating database")
	a.storage, err = clickhouse.New(a.config)
	if err != nil {
		return nil, fmt.Errorf("error while create creating database: %s", err.Error())
	}

	//appLog.Debug("Initializing database")
	//if err := a.storage.Init(); err != nil {
	//	return nil, fmt.Errorf("error while intializing clickhouse: %s", err.Error())
	//}

	appLog.Debug("Creating session caching")
	a.cache, err = RedisCache.New(a.config)
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

	a.logger.Info("Waiting application to exit")

	if err := a.srv.Network.Shutdown(ctx); err != nil {
		return fmt.Errorf("can't exit http server: %s", err.Error())
	}

	if err := a.cache.Client.Close(); err != nil {
		return fmt.Errorf("can't exit redis: %s", err.Error())
	}

	if err := a.storage.Db.Close(); err != nil {
		return fmt.Errorf("can't exit clickhouse: %s", err.Error())
	}

	return nil
}
