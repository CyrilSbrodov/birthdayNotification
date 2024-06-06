package handlers

import (
	"birthdayNotification/internal/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h *Handler) Subscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notification models.Notification
		notification.UserID = r.Context().Value(ctxKeyUser).(string)

		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if err := h.service.Subscribe(r.Context(), &notification); err != nil {
			if errors.Is(err, models.ErrorUserNotFound) {
				http.Error(w, "employee not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) Unsubscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var notification models.Notification
		notification.UserID = r.Context().Value(ctxKeyUser).(string)

		if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if err := h.service.Unsubscribe(r.Context(), &notification); err != nil {
			if errors.Is(err, models.ErrorUserNotFound) {
				http.Error(w, "employee not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) GetAllSubscribes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ctxKeyUser).(string)
		notifications, err := h.service.GetAllNotifications(r.Context(), id)
		if err != nil {
			if errors.Is(err, models.ErrorSubscribesNotFound) {
				http.Error(w, "subscribes not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(&notifications)
		if err != nil {
			http.Error(w, "failed to marshal to json", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}
