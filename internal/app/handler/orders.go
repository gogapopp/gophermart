package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gogapopp/gophermart/config"
	"github.com/gogapopp/gophermart/internal/app/models"
)

// userOrdersPostHandler загружает номер заказа пользователя для расчёта
func (h *Handler) userOrdersPostHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDkey).(int)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error read body", http.StatusInternalServerError)
		return
	}
	number, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(w, "incorrect request format", http.StatusBadRequest)
		return
	}
	if ok := Valid(number); !ok {
		http.Error(w, "unvalid order number", http.StatusUnprocessableEntity)
		return
	}
	order, err := OrderReq(number)
	if err != nil {
		if errors.As(err, &pgErr) {
			http.Error(w, "order already exists", http.StatusConflict)
			return
		} else if errors.Is(err, ErrTooManyRequests) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		} else {
			http.Error(w, "service reject", http.StatusInternalServerError)
			return
		}
	}
	Order := models.Order{
		Number:     number,
		Status:     order.Status,
		Accrual:    order.Accrual,
		UploadedAt: time.Now().Format(time.RFC3339),
	}
	orderID, err := h.services.Orders.Create(userID, Order)
	h.log.Infoln("POST /api/user/orders",
		fmt.Sprintf("Orders struct:"),
		fmt.Sprintf("Number %d", order.Number),
		fmt.Sprintf("Status %s", order.Status),
		fmt.Sprintf("Accrual %f", order.Accrual),
		fmt.Sprintf("orderID %d", orderID),
	)
	w.WriteHeader(http.StatusAccepted)
}

// userOrdersGetHandler  получает список загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *Handler) userOrdersGetHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GET /api/user/orders")
	userID := r.Context().Value(userIDkey).(int)

	orders, err := h.services.Orders.GetUserOrders(userID)
	if err != nil {
		h.log.Infoln("get error: ", err)
		http.Error(w, "error get user orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// Valid check number is valid or not based on Luhn algorithm
func Valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

type RespOrder struct {
	Number  int     `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func OrderReq(number int) (RespOrder, error) {
	var Order RespOrder
	client := resty.New()
	url := fmt.Sprintf("%s/api/orders/%d", config.AccSysAddr, number)

	resp, err := client.R().Get(url)
	if err != nil {
		return Order, err
	}

	if resp.StatusCode() == http.StatusOK {
		err = json.Unmarshal(resp.Body(), &Order)
		if err != nil {
			return Order, err
		}
	} else if resp.StatusCode() == http.StatusTooManyRequests {
		return Order, ErrTooManyRequests
	}
	return Order, err
}
