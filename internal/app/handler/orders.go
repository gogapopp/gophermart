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
	"github.com/gogapopp/gophermart/internal/app/storage"
)

// userOrdersPostHandler загружает номер заказа пользователя для расчёта
func (h *Handler) userOrdersPostHandler(w http.ResponseWriter, r *http.Request) {
	// получаем userID из контекста который был установлен мидлвеером userIdentity
	userID := r.Context().Value(userIDkey).(int)
	// читаем бади пост запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error read body", http.StatusInternalServerError)
		return
	}
	// прочитаный номер из запроса конвертируем в int
	number, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(w, "incorrect request format", http.StatusBadRequest)
		return
	}
	// проверяем валидность с помощью алгоритма Луна
	if ok := Valid(number); !ok {
		http.Error(w, "unvalid order number", http.StatusUnprocessableEntity)
		return
	}
	// отправляем запрос на сервер расчёта баллов и получаем структуру
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
	// создаём структуру для отправки в БД
	Order := models.Order{
		Number:     fmt.Sprintf("%d", number),
		Status:     order.Status,
		Accrual:    order.Accrual,
		UploadedAt: time.Now().Format(time.RFC3339),
	}
	// проверяем есть ли у юзера ордер с таким номером
	err = h.services.Orders.CheckUserOrder(userID, Order)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, storage.ErrUserRepeatValue) {
			http.Error(w, "user order already exist", http.StatusOK)
			return
		} else if errors.Is(err, storage.ErrRepeatValue) {
			http.Error(w, "order already exist", http.StatusConflict)
			return
		}
		http.Error(w, "error check user order", http.StatusInternalServerError)
		return
	}
	// если сервис вернул пустой order.status ставим ему значение "NEW"
	if len(Order.Status) < 1 {
		Order.Status = "NEW"
	}
	// создаём ордер в бд
	_, err = h.services.Orders.Create(userID, Order)
	if err != nil {
		if errors.Is(err, storage.ErrRepeatValue) {
			http.Error(w, "order already exist", http.StatusConflict)
			return
		}
		http.Error(w, "error create order", http.StatusInternalServerError)
		return
	}
	// обновляем баланс юзера (отправляем баллы лояльности которые нам расчитал сервис)
	err = h.services.Balance.UpdateUserBalance(userID, Order.Accrual)
	if err != nil {
		http.Error(w, "error check user order", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// userOrdersGetHandler  получает список загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *Handler) userOrdersGetHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("GET /api/user/orders")
	// получаем userID из контекста который был установлен мидлвеером userIdentity
	userID := r.Context().Value(userIDkey).(int)
	// получаем из БД все ордеры юзера
	orders, err := h.services.Orders.GetUserOrders(userID)
	if err != nil {
		http.Error(w, "error get user orders", http.StatusInternalServerError)
		return
	}
	// проверяем есть ли ордеры у юзера
	if len(orders) == 0 {
		http.Error(w, "users don't have orders", http.StatusNoContent)
		return
	}
	// возвращаем json с ордерами юзера
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	Number     string  `json:"order"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

// OrderReq отправляет запрос на сервис расчёта баллов
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
