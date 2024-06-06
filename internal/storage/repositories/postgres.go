package repositories

import (
	"birthdayNotification/cmd/loggers"
	"birthdayNotification/internal/config"
	"birthdayNotification/internal/models"
	"birthdayNotification/pkg/client/postgres"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type PGStore struct {
	client postgres.Client
	cfg    *config.Config
	logger *loggers.Logger
}

// createTable - функция создания новых таблиц в БД.
func createTable(ctx context.Context, client postgres.Client, logger *loggers.Logger) error {
	tx, err := client.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error("failed to begin transaction", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	//создание таблиц
	tables := []string{
		`CREATE TABLE IF NOT EXISTS employees (
            id SERIAL PRIMARY KEY,
            first_name VARCHAR(50) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
    		birthday DATE NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            login VARCHAR(50) NOT NULL UNIQUE,
            password_hash VARCHAR(255) NOT NULL,
            email VARCHAR(100) NOT NULL UNIQUE,
    		birthday DATE
        )`,
		`CREATE TABLE IF NOT EXISTS notifications (
    		id SERIAL PRIMARY KEY,
    		user_id INT NOT NULL,
    		employee_id INT NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
            FOREIGN KEY (employee_id) REFERENCES employees (id) ON DELETE CASCADE,
            UNIQUE (user_id, employee_id)
        )`,
	}

	for _, table := range tables {
		_, err = tx.Exec(ctx, table)
		if err != nil {
			logger.Error("Unable to create table", err)
			return err
		}
	}
	return tx.Commit(ctx)
}

func NewPGStore(client postgres.Client, cfg *config.Config, logger *loggers.Logger) (*PGStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := createTable(ctx, client, logger); err != nil {
		logger.Error("failed to create table", err)
		return nil, err
	}
	return &PGStore{
		client: client,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (p *PGStore) NewUser(ctx context.Context, u *models.User) error {
	hashPassword := p.hashPassword(u.Password)

	q := `INSERT INTO users (login, password_hash, email) VALUES ($1, $2, $3)`
	if _, err := p.client.Exec(ctx, q, u.Login, hashPassword, u.Email); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.ErrorUserConflict
		}
		p.logger.Error("Failure to insert object into table", err)
		return err
	}

	q = `INSERT INTO employees (first_name, last_name, birthday) VALUES ($1, $2, $3)`
	if _, err := p.client.Exec(ctx, q, u.FirstName, u.LastName, u.Birthday); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.ErrorUserConflict
		}
		p.logger.Error("Failure to insert object into table", err)
		return err
	}
	return nil

}

func (p *PGStore) Auth(ctx context.Context, u *models.User) (string, error) {
	hashPassword := p.hashPassword(u.Password)
	q := `SELECT id FROM users WHERE login=$1 AND password_hash=$2`
	if err := p.client.QueryRow(ctx, q, u.Login, hashPassword).Scan(&u.Id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", models.ErrorUserNotFound
		}
		p.logger.Error("Failure to select object from table", err)
		return "", err
	}
	return u.Id, nil
}

func (p *PGStore) GetAllEmployees(ctx context.Context) ([]models.Employee, error) {
	var employees []models.Employee
	q := `SELECT id, first_name, last_name, birthday FROM employees`
	rows, err := p.client.Query(ctx, q)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrorUserNotFound
		}
		p.logger.Error("failure to select object from table")
	}
	for rows.Next() {
		var employee models.Employee
		if err = rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.Birthday); err != nil {
			p.logger.Error("failure to scan object from table")
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}
func (p *PGStore) Subscribe(ctx context.Context, n *models.Notification) error {
	q := `INSERT INTO notifications (user_id, employee_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	if _, err := p.client.Exec(ctx, q, n.UserID, n.EmployeeID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); err != nil && pgErr.Code == "23503" {
			return models.ErrorUserNotFound
		}
		p.logger.Error("failure to insert object into table")
		return err
	}
	return nil
}

func (p *PGStore) Unsubscribe(ctx context.Context, n *models.Notification) error {
	q := `DELETE FROM notifications WHERE user_id=$1 AND employee_id=$2`
	if _, err := p.client.Exec(ctx, q, n.UserID, n.EmployeeID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); err != nil && pgErr.Code == "23503" {
			return models.ErrorUserNotFound
		}
		p.logger.Error("failure to delete object from table")
		return err
	}
	return nil
}

func (p *PGStore) GetAllSubscribes(ctx context.Context, id string) ([]models.Notification, error) {
	q := `SELECT id, user_id, employee_id FROM notifications WHERE user_id=$1`
	rows, err := p.client.Query(ctx, q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrorSubscribesNotFound
		}
		p.logger.Error("failure to select object from table")
		return nil, err
	}
	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		if err = rows.Scan(&notification.ID, &notification.UserID, &notification.EmployeeID); err != nil {
			p.logger.Error("failure to scan object from table")
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (p *PGStore) Notifications(ctx context.Context, id string) ([]models.Employee, error) {
	q := `SELECT e.id AS employee_id, e.first_name, e.last_name
	FROM employees e JOIN notifications n ON e.id = n.employee_id
	WHERE n.user_id = $1
	AND DATE_PART('day', e.birthday) = DATE_PART('day', CURRENT_DATE)
	AND DATE_PART('month', e.birthday) = DATE_PART('month', CURRENT_DATE)`

	rows, err := p.client.Query(ctx, q, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrorSubscribesNotFound
		}
		p.logger.Error("failure to select object from table")
		return nil, err
	}
	defer rows.Close()
	var employees []models.Employee
	for rows.Next() {
		var employee models.Employee
		if err = rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName); err != nil {
			p.logger.Error("failure to scan object from table")
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (p *PGStore) AddEmployees(ctx context.Context, employees *[]models.Employee) error {
	q := `INSERT INTO employees (first_name, last_name, birthday) VALUES ($1, $2, $3)`
	batch := &pgx.Batch{}
	for _, employee := range *employees {
		batch.Queue(q, employee.FirstName, employee.LastName, employee.Birthday)
	}
	br := p.client.SendBatch(ctx, batch)
	if _, err := br.Exec(); err != nil {
		p.logger.Error("Failure to insert object into table", err)
		return err
	}
	return nil
}

func (p *PGStore) hashPassword(pass string) string {
	h := hmac.New(sha256.New, []byte("password"))
	h.Write([]byte(pass))
	return fmt.Sprintf("%x", h.Sum(nil))
}
