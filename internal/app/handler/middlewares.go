package handler

import (
	"context"
	"net/http"
	"strings"
)

func (h *Handler) userIdentity(v http.HandlerFunc) http.HandlerFunc {
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

		h.ctx = context.WithValue(h.ctx, "userID", userID)
		v.ServeHTTP(w, r)
	})
}
