package interfaces

import "net/http"

type Handlers interface {
	SearchUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type Storer interface {
	SearchUsers() ([]byte, error)
	CreateUser(name, email string) (string, error)
	GetUser(id string) ([]byte, error)
	UpdateUser(id, name, email string) error
	DeleteUser(id string) error
}
