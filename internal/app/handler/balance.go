package handler

import (
	"encoding/json"
	"net/http"
)

// userBalanceGetHanlder возвращает структуру баланса юзера
func (h *Handler) userBalanceGetHanlder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GET /api/user/balance")
	// получаем userID из контекста который был установлен мидлвеером userIdentity
	userID := r.Context().Value(userIDkey).(int)
	// получаем юзер баланс из БД
	userBalance, err := h.services.GetUserBalance(userID)
	if err != nil {
		http.Error(w, "error get user balance", http.StatusInternalServerError)
		return
	}
	// возвращаем баланс юзера в виде json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userBalance); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
