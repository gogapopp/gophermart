package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gogapopp/gophermart/internal/app/service"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

var pgErr *pgconn.PgError

type Handler struct {
	services *service.Service
	log      *zap.SugaredLogger
}

func NewHandler(services *service.Service, log *zap.SugaredLogger) *Handler {
	return &Handler{services: services, log: log}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/user/register", h.userRegisterPostHandler)
	r.Post("/api/user/login", h.userLoginPostHandler)
	r.Post("/api/user/orders", h.userIdentity(h.userOrdersPostHandler))
	r.Get("/api/orders/{numbers}", (h.GetOrder))
	r.Get("/api/user/orders", h.userIdentity(h.userOrdersGetHandler))
	r.Get("/api/user/balance", h.userIdentity(h.userBalanceGetHanlder))
	r.Post("/api/user/balance/withdraw", h.userIdentity(h.userBalanceWithdrawPostHanlder))
	r.Get("/api/user/withdrawals", h.userIdentity(h.userBalanceWithdrawalsGetHanlder))
	r.NotFound(http.NotFound)

	return r
}
