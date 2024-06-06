package service

import (
	"birthdayNotification/internal/models"
	"context"
)

type Service interface {
	CreateUser(ctx context.Context, u *models.User) error
	GenerateToken(ctx context.Context, u *models.User) (string, error)
	ParseToken(ctx context.Context, accessToken string) (string, error)
	CollectEmployeesFromAPI(ctx context.Context, employees *[]models.Employee) error
	GetAll(ctx context.Context) ([]models.Employee, error)
	GetAllNotifications(ctx context.Context, id string) ([]models.Notification, error)
	Subscribe(ctx context.Context, n *models.Notification) error
	Unsubscribe(ctx context.Context, n *models.Notification) error
}
