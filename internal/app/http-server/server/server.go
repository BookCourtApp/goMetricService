package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/config"
)

type Server struct {
	Network *http.Server
}

func New(config *config.Config, router *chi.Mux) *Server {

	httpSrv := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	return &Server{
		Network: httpSrv,
	}
}
