package handlers

import (
	"birthdayNotification/cmd/loggers"
	"birthdayNotification/internal/config"
	"birthdayNotification/internal/service"
	"github.com/gorilla/mux"
)

type Handlers interface {
	Register(router *mux.Router)
}

type Handler struct {
	cfg     *config.Config
	logger  *loggers.Logger
	service service.Service
}

func NewHandler(cfg *config.Config, logger *loggers.Logger, service service.Service) *Handler {
	return &Handler{
		cfg:     cfg,
		logger:  logger,
		service: service,
	}
}

func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/api/employees", h.AddEmployees()).Methods("POST")
	r.HandleFunc("/api/employees", h.GetAllEmployees()).Methods("GET")
	r.HandleFunc("/api/register", h.SignUp()).Methods("POST")
	r.HandleFunc("/api/login", h.SignIn()).Methods("POST")
	secure := r.PathPrefix("/auth").Subrouter()
	secure.Use(h.userIdentity)
	secure.HandleFunc("/api/notification/subscribe", h.Subscribe()).Methods("POST")
	secure.HandleFunc("/api/notification/unsubscribe", h.Unsubscribe()).Methods("POST")
	secure.HandleFunc("/api/notification", h.GetAllSubscribes()).Methods("GET")
}
