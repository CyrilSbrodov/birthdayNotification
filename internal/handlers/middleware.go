package handlers

import (
	"context"
	"net/http"
	"strings"
)

const (
	autharizationHeader = "Authorization"
	ctxKeyUser          = "userId"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(autharizationHeader)
		if header == "" {
			http.Error(w, "empty auth header", http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		userId, err := h.service.ParseToken(r.Context(), headerParts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, userId)))
	})
}
