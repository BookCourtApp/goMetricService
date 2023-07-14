package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/config"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/http-server/handlers/sendData"
)

type App struct {
	config *config.Config
	srv    http.Server
}

func New() (*App, error) {
	a := &App{}

	a.config = config.MustLoad()

	router := chi.NewRouter()
	router.Get("/test", sendData.New())

	a.srv = http.Server{
		Addr:    a.config.Address,
		Handler: router,
	}

	return a, nil
}

func (a *App) Run() error {
	log.Printf("Listening address: %s", a.srv.Addr)

	if err := a.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("error occured while listening port: %s", err.Error())
	}

	return nil
}
