package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gogapopp/gophermart/config"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	log        *zap.SugaredLogger
}

func NewServer(log *zap.SugaredLogger) *Server {
	return &Server{httpServer: &http.Server{}, log: log}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           config.RunAddr,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	s.log.Info("run server at address: ", config.RunAddr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("server shutdown...")
	return s.httpServer.Shutdown(ctx)
}
