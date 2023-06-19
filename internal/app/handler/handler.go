package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gogapopp/gophermart/internal/app/service"
	"go.uber.org/zap"
)

type Handler struct {
	ctx      context.Context
	services *service.Service
	log      *zap.SugaredLogger
}

func NewHandler(ctx context.Context, services *service.Service, log *zap.SugaredLogger) *Handler {
	return &Handler{ctx: ctx, services: services, log: log}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/user/register", h.userRegisterHandler)
	r.Post("/api/user/login", h.userLoginHandler)
	r.Post("/api/user/orders", h.userIdentity(h.userOrdersPostHandler))
	r.Get("/api/user/orders", h.userOrdersGetHandler)
	r.Get("/api/user/balance", h.userBalanceGetHanlder)
	r.Post("/api/user/balance/withdraw", h.userBalanceWithdrawHanlder)
	r.Get("/api/user/withdrawals", h.userBalanceWithdrawalsHanlder)
	r.NotFound(http.NotFound)

	return r
}
