package models

import "errors"

var (
	ErrorUserConflict       = errors.New("user or email already exists")
	ErrorUserNotFound       = errors.New("user not found")
	ErrorSubscribesNotFound = errors.New("subscribes not found")
)
