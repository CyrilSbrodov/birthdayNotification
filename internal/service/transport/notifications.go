package transport

import (
	"birthdayNotification/internal/models"
	"context"
)

func (t *Transport) CollectEmployeesFromAPI(ctx context.Context, employees *[]models.Employee) error {
	return t.storage.AddEmployees(ctx, employees)
}

func (t *Transport) GetAll(ctx context.Context) ([]models.Employee, error) {
	return t.storage.GetAllEmployees(ctx)
}

func (t *Transport) GetAllNotifications(ctx context.Context, id string) ([]models.Notification, error) {
	return t.storage.GetAllSubscribes(ctx, id)
}

func (t *Transport) Subscribe(ctx context.Context, n *models.Notification) error {
	if err := t.storage.Subscribe(ctx, n); err != nil {
		return err
	}

	subscriptions, err := t.storage.GetAllSubscribes(ctx, n.UserID)
	if err != nil {
		return err
	}

	if len(subscriptions) == 1 { // Первая подписка
		t.notification.Start(n.UserID)
	}
	return nil
}

func (t *Transport) Unsubscribe(ctx context.Context, n *models.Notification) error {
	if err := t.storage.Unsubscribe(ctx, n); err != nil {
		return err
	}

	subscriptions, err := t.storage.GetAllSubscribes(ctx, n.UserID)
	if err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		t.notification.Stop()
	}

	return nil
}

func (t *Transport) GetTodayBirthday(ctx context.Context, id string) ([]models.Employee, error) {
	return t.storage.Notifications(ctx, id)
}
