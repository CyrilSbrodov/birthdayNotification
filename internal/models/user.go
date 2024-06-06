package models

import "time"

type User struct {
	Id        string    `json:"id"`
	Login     string    `json:"login"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Birthday  time.Time `json:"birthday"`
}
