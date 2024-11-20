package server

import (
	"context"
	"net/http"
	"todo-app/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func CreateNewServer(config *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: ":" + config.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    config.HTTP.ReadTimeout,
			WriteTimeout:   config.HTTP.WriteTimeout,
			MaxHeaderBytes: config.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

