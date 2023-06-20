package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

// userRegisterHandler регистрирует пользователя
func (h *Handler) userRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("POST /api/user/register")
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	_, err := h.services.Auth.CreateUser(req)
	if err != nil {
		if errors.As(err, &pgErr) {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		http.Error(w, "error create user", http.StatusInternalServerError)
		return
	}

	token, err := h.services.Auth.GenerateToken(req.Login, req.Password)
	if err != nil {
		http.Error(w, "error generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// userLoginHandler аутентифицирует пользователя
func (h *Handler) userLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("POST /api/user/login")
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	token, err := h.services.Auth.GenerateToken(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrNoRows) {
			http.Error(w, "wrong login or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "error generate token", http.StatusInternalServerError)
		return
	}

	h.log.Info(fmt.Sprintf("Bearer %s", token))
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
