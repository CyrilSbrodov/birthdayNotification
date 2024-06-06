package handlers

import (
	"birthdayNotification/internal/models"
	"encoding/json"
	"errors"
	"net/http"
)

func (h *Handler) AddEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employee []models.Employee
		if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := h.service.CollectEmployeesFromAPI(r.Context(), &employee); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) GetAllEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employees, err := h.service.GetAll(r.Context())
		if err != nil {
			if errors.Is(err, models.ErrorUserNotFound) {
				http.Error(w, "Users not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(employees)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "failed to encode data", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}
