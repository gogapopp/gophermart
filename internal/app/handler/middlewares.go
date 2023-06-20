package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
)

type ctxKey string

var userIDkey ctxKey = "userID"

func (h *Handler) userIdentity(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "empty authorization header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		userID, err := h.services.Auth.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), userIDkey, userID))
		next(w, r)
	})
}

// RequestLogger логирует POST запрос
func (h *Handler) RequestLogger(v http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// читаем боди запоса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Info(err)
		}
		// возвращаем данные обратно
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		v(w, r)
		h.log.Info("POST request",
			r.Method,
			string(body),
		)
	})
}
