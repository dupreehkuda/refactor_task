package handlers

import (
	"net/http"
)

type UserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *UserRequest) Bind(r *http.Request) error { return nil }
