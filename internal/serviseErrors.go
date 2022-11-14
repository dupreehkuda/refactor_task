package internal

import "errors"

var (
	UserNotFound = errors.New("user_not_found")
	UserExists   = errors.New("user_exists")
)
