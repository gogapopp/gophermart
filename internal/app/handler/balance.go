package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) userBalanceGetHanlder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GET /api/user/balance")
	userID := r.Context().Value(userIDkey).(int)
	userBalance, err := h.services.GetUserBalance(userID)
	if err != nil {
		http.Error(w, "error get user balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userBalance); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
