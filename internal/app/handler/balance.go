package handler

import (
	"encoding/json"
	"net/http"
)

// getUserBalanceHanlder возвращает структуру баланса юзера
func (h *Handler) getUserBalanceHanlder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// получаем userID из контекста который был установлен мидлвеером userIdentity
	userID := ctx.Value(userIDkey).(int)
	// получаем юзер баланс из БД
	userBalance, err := h.services.GetUserBalance(ctx, userID)
	if err != nil {
		http.Error(w, ErrGetBalance.Error(), http.StatusInternalServerError)
		return
	}
	// возвращаем баланс юзера в виде json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userBalance); err != nil {
		http.Error(w, ErrEncogingResp.Error(), http.StatusInternalServerError)
		return
	}
}
