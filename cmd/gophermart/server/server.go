package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gogapopp/gophermart/config"
	"go.uber.org/zap"
)

// А почему решил хранить структуру и инициализацию Server в директории `cmd`? Тут логичнее смотрится работа самого main и всяких config/init
// Это не ошибка, тут на свое усмотрение, но просто интересна логика
type Server struct {
	httpServer *http.Server
	log        *zap.SugaredLogger
}

// NewServer возвращает структуру сервера
func NewServer(log *zap.SugaredLogger) *Server {
	return &Server{httpServer: &http.Server{}, log: log}
}

// Run запускает сервер
func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         config.RunAddr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.log.Info("run server at address: ", config.RunAddr)
	return s.httpServer.ListenAndServe()
}

// Shutdown выключает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("server shutdown...")
	return s.httpServer.Shutdown(ctx)
}
