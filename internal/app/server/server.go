package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gogapopp/gophermart/internal/app/config"
	"go.uber.org/zap"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server
	log        *zap.SugaredLogger
}

// NewServer возвращает структуру сервера
func NewServer(config *config.Config, log *zap.SugaredLogger) *Server {
	return &Server{config: config, httpServer: &http.Server{}, log: log}
}

// Run запускает сервер
func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           s.config.RunAddr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.log.Info("run server at address: ", s.config.RunAddr)
	return s.httpServer.ListenAndServe()
}

// Shutdown выключает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("server shutdown...")
	return s.httpServer.Shutdown(ctx)
}
