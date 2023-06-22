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
	// получаем userID из контекста который был установлен мидлвеером userIdentity
	userID := r.Context().Value(userIDkey).(int)
	// получаем баланс юзера
	userBalance, err := h.services.GetUserBalance(userID)
	if err != nil {
		http.Error(w, "error get user balance", http.StatusInternalServerError)
		return
	}
	// читаем бади пост запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error read request body", http.StatusInternalServerError)
		return
	}
	// декодируем бади в структуру
	var rb models.RequestBody
	err = json.Unmarshal(body, &rb)
	if err != nil {
		http.Error(w, "error unmarshal request body", http.StatusInternalServerError)
		return
	}
	// конвертируем номер в int
	number, err := strconv.Atoi(rb.Order)
	if err != nil {
		http.Error(w, "incorrect request format", http.StatusBadRequest)
		return
	}
	// проверяем номер на валидность
	if ok := Valid(number); !ok {
		http.Error(w, "unvalid order number", http.StatusUnprocessableEntity)
		return
	}
	// проверяем достаточно ли денег у юзера на балансе (rb.Sum это сумма списания баллов которую мы получили из бади)
	if userBalance.Current <= rb.Sum {
		http.Error(w, "payment required", http.StatusPaymentRequired)
		return
	}
	// создаём структуру для записи withdraw в БД
	WithdrawResp := models.Withdraw{
		UserID:      userID,
		Order:       rb.Order,
		Sum:         rb.Sum,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}
	// обновляем баланс юзера (вычитаем из баланса сумму списания баллов)
	err = h.services.Balance.UpdateUserBalance(userID, -rb.Sum)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error update user balance", http.StatusInternalServerError)
		return
	}
	// записываем withdraw юзера в БД
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
	// получаем userID из контекста который был установлен мидлвеером userIdentity
	userID := r.Context().Value(userIDkey).(int)
	// получаем все withdraw юзера из БД
	userWithdrawals, err := h.services.GetUserWithdrawals(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error get user withdrawals", http.StatusInternalServerError)
		return
	}
	// проверяем есть ли у юзера вообще withdraw
	if len(userWithdrawals) == 0 {
		http.Error(w, "users don't have withdrawals", http.StatusNoContent)
		return
	}
	// выводим все withdraw юзера в виде json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userWithdrawals); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
