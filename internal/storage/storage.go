package storage

import (
	"birthdayNotification/internal/models"
	"context"
)

type Storage interface {
	NewUser(ctx context.Context, u *models.User) error
	Auth(ctx context.Context, u *models.User) (string, error)
	AddEmployees(ctx context.Context, employees *[]models.Employee) error
	GetAllEmployees(ctx context.Context) ([]models.Employee, error)
	Subscribe(ctx context.Context, n *models.Notification) error
	Unsubscribe(ctx context.Context, n *models.Notification) error
	GetAllSubscribes(ctx context.Context, id string) ([]models.Notification, error)
	Notifications(ctx context.Context, id string) ([]models.Employee, error)
}
