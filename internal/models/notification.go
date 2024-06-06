package models

type Notification struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	EmployeeID string `json:"employee_id"`
}
