package storage

import "time"

type User struct {
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email,omitempty"`
}

type UserList map[string]User

type UserStore struct {
	Increment int64    `json:"increment"`
	List      UserList `json:"list"`
}
