package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gogapopp/gophermart/internal/app/config"
	"github.com/gogapopp/gophermart/internal/app/logger"
	"github.com/gogapopp/gophermart/internal/app/service"
	"go.uber.org/zap"
)

type Handler struct {
	config   *config.Config
	services *service.Service
	log      *zap.SugaredLogger
}

// NewHandler возвращает экземляр структуры хендлера
func NewHandler(config *config.Config, services *service.Service, log *zap.SugaredLogger) *Handler {
	return &Handler{config: config, services: services, log: log}
}

// InitRoutes инициализирует роуты
func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", logger.ResponseLogger(h.postUserRegisterHandler))
		r.Post("/login", logger.ResponseLogger(h.postUserLoginHandler))
		r.Post("/orders", logger.ResponseLogger(h.userIdentity(h.postUserOrdersHandler)))
		r.Get("/orders", logger.RequestLogger(h.userIdentity(h.getUserOrdersHandler)))
		r.Get("/balance", logger.RequestLogger(h.userIdentity(h.getUserBalanceHanlder)))
		r.Post("/balance/withdraw", logger.ResponseLogger(h.userIdentity(h.postUserBalanceWithdrawHanlder)))
		r.Get("/withdrawals", logger.RequestLogger(h.userIdentity(h.getUserBalanceWithdrawalsHanlder)))
	})
	r.NotFound(http.NotFound)

	return r
}
