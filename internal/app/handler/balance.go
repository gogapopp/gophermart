package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) userBalanceGetHanlder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GET /api/user/balance")
	userID := r.Context().Value(userIDkey).(int)
	userBalance, err := h.services.GetUserBalance(userID)
	fmt.Println(userBalance)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userBalance); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) userBalanceWithdrawPostHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "balance withdraw")
}

func (h *Handler) userBalanceWithdrawalsGetHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "user balance withdrawals")
}
