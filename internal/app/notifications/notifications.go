package notifications

import (
	"birthdayNotification/internal/storage"
	"context"
	"fmt"
	"log"
	"time"
)

type Notifications struct {
	store  storage.Storage
	ticker *time.Ticker
	quit   chan struct{}
}

func NewNotifications(store storage.Storage) *Notifications {
	return &Notifications{
		store: store,
	}
}

func (n *Notifications) Start(id string) {
	n.ticker = time.NewTicker(24 * time.Hour)
	go func(id string) {
		for {
			select {
			case <-n.ticker.C:
				n.CheckBirthdays(id)
			case <-n.quit:
				n.ticker.Stop()
				return
			}
		}
	}(id)
}

func (n *Notifications) Stop() {
	close(n.quit)
}

func (n *Notifications) CheckBirthdays(id string) {
	ctx := context.Background()
	notifications, err := n.store.Notifications(ctx, id)
	if err != nil {
		log.Println("Error checking birthdays:", err)
		return
	}

	for _, notification := range notifications {
		// Отправка email-уведомленияне не реализована
		fmt.Println("Сегодня День Рождения у ", notification.FirstName, notification.LastName)
	}
}
