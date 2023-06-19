package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) userBalanceGetHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "user balance")
}

func (h *Handler) userBalanceWithdrawHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "balance withdraw")
}

func (h *Handler) userBalanceWithdrawalsHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "user balance withdrawals")
}
