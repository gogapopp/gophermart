package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gogapopp/gophermart/internal/app/models"
)

func (h *Handler) userBalanceWithdrawPostHanlder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("POST /api/user/balance/withdraw")
	userID := r.Context().Value(userIDkey).(int)
	userBalance, err := h.services.GetUserBalance(userID)
	if err != nil {
		http.Error(w, "error get user balance", http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error read request body", http.StatusInternalServerError)
		return
	}
	var rb models.RequestBody
	err = json.Unmarshal(body, &rb)
	if err != nil {
		http.Error(w, "error unmarshal request body", http.StatusInternalServerError)
		return
	}
	number, err := strconv.Atoi(rb.Order)
	if err != nil {
		http.Error(w, "incorrect request format", http.StatusBadRequest)
		return
	}
	if ok := Valid(number); !ok {
		http.Error(w, "unvalid order number", http.StatusUnprocessableEntity)
		return
	}
	if userBalance.Current <= rb.Sum {
		http.Error(w, "payment required", http.StatusPaymentRequired)
		return
	}
	WithdrawResp := models.Withdraw{
		UserID:      userID,
		Order:       rb.Order,
		Sum:         rb.Sum,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	err = h.services.Withdrawals.UserWithdraw(userID, WithdrawResp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error response withdraw", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) userBalanceWithdrawalsGetHanlder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GET /api/user/balance/withdrawals")
	userID := r.Context().Value(userIDkey).(int)
	userWithdrawals, err := h.services.GetUserWithdrawals(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error get user withdrawals", http.StatusInternalServerError)
		return
	}
	if len(userWithdrawals) == 0 {
		http.Error(w, "users don't have withdrawals", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userWithdrawals); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
