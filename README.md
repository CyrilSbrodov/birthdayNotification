#### Оглавление:
____
0. [Сервис уведомления о Днях Рождения](#сервис-уведомления-о-Днях-Рождения).
1. [ЗАВИСИМОСТИ](#зависимости).
2. [ЗАПУСК/СБОРКА](#запусксборка).
2.1. [Конфигурация](#конфигурация).
2.1.1 [Конфигурационный файл](#3-конфигурационный-файл).
2.2. [Запуск сервера](#запуск-сервера).
3. [Для разработчиков](#для-разработчиков).
3.1. [Сервер](#2-сервер).
____

# Сервис уведомления о Днях Рождения.

Данный сервис позволяет загрузить в базу список человек(сотрудников) с фамилиями, именами и датой рождения. Реализована регистрация и аутентификация. Реализованы подписки на Дни Рождения, чтобы получить уведомление о наступившем дне.
Возможно реализовать отправку уведомлений через ТГ бота или почту. Возможно реализовать детальный вывод всех сотрудников списком.

Структура [сотрудников](https://github.com/CyrilSbrodov/birthdayNotification/blob/main/internal/models/employee.go):
```GO
type Employee struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
}
```
Структура [пользователей](https://github.com/CyrilSbrodov/birthdayNotification/blob/main/internal/models/user.go):
```GO
type User struct {
	Id        string    `json:"id"`
	Login     string    `json:"login"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Birthday  time.Time `json:"birthday"`
}
```
Структура [уведомлений](https://github.com/CyrilSbrodov/birthdayNotification/blob/main/internal/models/notification.go)
```GO
type Notification struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	EmployeeID string `json:"employee_id"`
}
```

Структура сервиса следующая:
1) Хэндлеры - обработка полученных данных и отправка их в транспортный узел.
2) Транспорт - обработка данных и отправка их в слой репозитория.
3) Слой БД - взаимодействие с БД PostgreSQL
4) Слой уведомлений - проверка наличия подписок и уведомлений.
____
# ЗАВИСИМОСТИ.

Используется язык go версии 1.22. Используемые библиотеки:
- github.com/BurntSushi/toml v1.4.0
- github.com/dgrijalva/jwt-go v3.2.0+incompatible
- github.com/gorilla/mux v1.8.1
- github.com/ilyakaznacheev/cleanenv v1.5.0
- github.com/jackc/pgpassfile v1.0.0
- github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9
- github.com/jackc/pgx/v5 v5.6.0
- github.com/jackc/puddle/v2 v2.2.1t
- github.com/joho/godotenv v1.5.1
- golang.org/x/crypto v0.24.0
- golang.org/x/sync v0.7.0
- golang.org/x/text v0.16.0
- gopkg.in/yaml.v3 v3.0.1
- olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3
- POSTGRESQL latest
____

# ЗАПУСК/СБОРКА

## Конфигурация

1) конфигурационный файл

```
---
env: "local" # local, dev, prod
storage_path: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
listener:
  addr: "localhost:8082"
  timeout: 4s
  idle_timeout: 60s
```

Возможность выбора адреса и порта сервера. Выбрать таймауты. И виды логов в зависимости от env.

## Запуск сервера

Необходимо запустить сервер из пакета [cmd](https://github.com/CyrilSbrodov/birthdayNotification/blob/main/cmd/main.go)
```
cd cmd
go run main.go //go run cmd/main.go
```
Или запустить через команду make:
```
make run
```
 
# Для разработчиков
Структура приложения позволяет нативно вносить корректировки:
[Структура сервера](https://github.com/CyrilSbrodov/birthdayNotification/blob/main/internal/app/app.go):
```GO
type ServerApp struct {
	cfg    config.Config //конфиг
	logger *loggers.Logger //логгер
	router *mux.Router //роутер
}
```

Немного о сервере:

Сервер получает данные по следующим эндпоинтам:
1) [http](https://github.com/CyrilSbrodov/birthdayNotification/blob/main/internal/handlers/handler.go):
```GO
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
```

Дополнительно можно реализовать тестирование, взаимодействие с фронтом.

Через docker compose up -d можно заупстить БД Postgres.
Формат приема json:
```
"/api/employees" списком:
[
    {
        "first_name":"Test",
        "last_name":"TestTest",
        "birthday":"2001-06-06T00:00:00Z"
    },
{
        "first_name":"Test1",
        "last_name":"TestTest1",
        "birthday":"2020-06-06T00:00:00Z"
    }
]

"/auth/api/notification/subscribe":
{
    "employee_id":"4"
}
```

Таблицы БД создаются автоматически при первом запуске, если таких таблиц нет. Миграции делать не стал.
```
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
```
