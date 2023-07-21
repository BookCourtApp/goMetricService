package app

import (
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/config"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/http-server/handler"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/http-server/server"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/service"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage/clickhouse"
)

type App struct {
	config  *config.Config
	srv     *server.Server
	storage *clickhouse.Clickhouse
}

func New() (*App, error) {
	a := &App{}

	a.config = config.MustLoad()

	db, err := clickhouse.New()
	if err != nil {
		return nil, err
	}
	a.storage = db

	a.storage.Test() // тест загрузки метрики в бд, потом убрать

	srvc := service.New()

	router := chi.NewRouter()
	router.Get("/test", handler.SaveHandler(srvc))

	a.srv = server.New(a.config, router)

	return a, nil
}

func (a *App) Run() error {
	log.Printf("Listening address: %s", a.srv.Network.Addr)

	if err := a.srv.Network.ListenAndServe(); err != nil {
		return fmt.Errorf("error occured while listening port: %s", err.Error())
	}

	return nil
}
