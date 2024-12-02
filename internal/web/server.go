package web

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Negat1v9/work-marketplace/internal/config"
	"github.com/Negat1v9/work-marketplace/internal/services"
	"github.com/Negat1v9/work-marketplace/internal/storage"
)

type Server struct {
	Server   http.Server
	cfg      *config.Config
	log      *slog.Logger
	services *services.Services
	store    storage.Store
}

func NewServer(cfg *config.Config, log *slog.Logger, services *services.Services, store storage.Store) *Server {
	s := &Server{
		Server: http.Server{
			Addr: cfg.WebConfig.Port,
		},
		cfg:      cfg,
		log:      log,
		services: services,
		store:    store,
	}

	return s
}

func (s *Server) Run() error {

	s.InitRoutes()

	s.log.Info("start listen...")

	return s.Server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
