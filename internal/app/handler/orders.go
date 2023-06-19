package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/gogapopp/gophermart/config"
	"github.com/gogapopp/gophermart/models"
)

// userOrdersPostHandler загружает номер заказа пользователя для расчёта
func (h *Handler) userOrdersPostHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("userOrdersPostHandler called")
	userID := h.ctx.Value("userID")

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
	err = OrderReq(number)
	if err != nil {
		fmt.Println(err)
	}
	h.log.Info(Order)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "download order, userID: %d: %d", userID, number)
}

// userOrdersGetHandler  получает список загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
func (h *Handler) userOrdersGetHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("userOrdersGetHandler called")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "orders list")
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

var Order models.Order

func OrderReq(number int) error {
	client := resty.New()

	url := fmt.Sprintf("%s/api/orders/%d", config.AccSysAddr, number)

	resp, err := client.R().Get(url)
	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		err = json.Unmarshal(resp.Body(), &Order)
		if err != nil {
			return err
		}
		fmt.Printf("Order: %d\n", Order.Number)
		fmt.Printf("Status: %s\n", Order.Status)
		fmt.Printf("Accrual: %d\n", Order.Accrual)
		fmt.Printf("Status: %s\n", Order.Status)
		fmt.Printf("Error: %s\n", resp.Status())
	} else {
		fmt.Printf("Error: %s\n", resp.Status())
	}
	return nil
}
