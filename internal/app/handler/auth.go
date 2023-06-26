package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

// postUserRegisterHandler регистрирует пользователя
func (h *Handler) postUserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// декодируем боди пост запроса
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrDecodingReq.Error(), http.StatusBadRequest)
		return
	}
	// отправляем запрос в бд на создание юзера
	_, err := h.services.Auth.CreateUser(req)
	if err != nil {
		if errors.As(err, &pgErr) {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		http.Error(w, "error create user", http.StatusInternalServerError)
		return
	}
	// получает jwt токен
	token, err := h.services.Auth.GenerateToken(req.Login, req.Password)
	if err != nil {
		http.Error(w, ErrGenerateToken.Error(), http.StatusInternalServerError)
		return
	}
	// записываем jwt токен в http заголовок
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

// postUserLoginHandler аутентифицирует пользователя
func (h *Handler) postUserLoginHandler(w http.ResponseWriter, r *http.Request) {
	// декодируем боди пост запроса
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrDecodingReq.Error(), http.StatusBadRequest)
		return
	}
	// получаем jwt токен (внутри GenerateToken шлём запрос на получение информации о юзере)
	token, err := h.services.Auth.GenerateToken(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			http.Error(w, "wrong login or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, ErrGenerateToken.Error(), http.StatusInternalServerError)
		return
	}
	// записываем jwt токен в http заголовок
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
