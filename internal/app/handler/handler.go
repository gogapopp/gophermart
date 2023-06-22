package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gogapopp/gophermart/internal/app/service"
	"go.uber.org/zap"
)

type Handler struct {
	services *service.Service
	log      *zap.SugaredLogger
}

// NewHandler возвращает экземляр структуры хендлера
func NewHandler(services *service.Service, log *zap.SugaredLogger) *Handler {
	return &Handler{services: services, log: log}
}

// InitRoutes инициализирует роуты
func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/user/register", h.userRegisterPostHandler)
	r.Post("/api/user/login", h.userLoginPostHandler)
	r.Post("/api/user/orders", h.userIdentity(h.userOrdersPostHandler))
	r.Get("/api/user/orders", h.userIdentity(h.userOrdersGetHandler))
	r.Get("/api/user/balance", h.userIdentity(h.userBalanceGetHanlder))
	r.Post("/api/user/balance/withdraw", h.userIdentity(h.userBalanceWithdrawPostHanlder))
	r.Get("/api/user/withdrawals", h.userIdentity(h.userBalanceWithdrawalsGetHanlder))
	r.NotFound(http.NotFound)

	return r
}
