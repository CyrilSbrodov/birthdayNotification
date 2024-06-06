package transport

import (
	"birthdayNotification/internal/app/notifications"
	"birthdayNotification/internal/storage"
	"birthdayNotification/internal/storage/repositories"
)

type Transport struct {
	storage      storage.Storage
	notification *notifications.Notifications
}

func NewTransport(repo repositories.PGStore, n *notifications.Notifications) *Transport {
	return &Transport{
		&repo,
		n,
	}
}
