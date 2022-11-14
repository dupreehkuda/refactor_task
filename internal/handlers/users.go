package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	i "github.com/dupreehkuda/refactor_task/internal"
)

// SearchUsers returns list of all users
func (h handlers) SearchUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := h.storage.SearchUsers()
	if err != nil {
		log.Printf("Error occurred getting list: %v", err)
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// CreateUser passes data to storage and creates new user
func (h handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	request := UserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id, err := h.storage.CreateUser(request.DisplayName, request.Email)

	switch err {
	case i.UserExists:
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	case nil:
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]string{
			"user_id": id,
		})
		return
	default:
		log.Printf("Error occurred creating user: %v", err)
		_ = render.Render(w, r, ErrServerError(err))
		return
	}
}

// GetUser passes id to storage to get specific user
func (h handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	resp, err := h.storage.GetUser(id)

	switch err {
	case i.UserNotFound:
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	default:
		log.Printf("Error occurred getting user: %v", err)
		_ = render.Render(w, r, ErrServerError(err))
		return
	}
}

// UpdateUser passes new user data to storage to update data
func (h handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	request := UserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id := chi.URLParam(r, "id")

	err := h.storage.UpdateUser(id, request.DisplayName, request.Email)

	switch err {
	case i.UserNotFound:
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	case nil:
		render.Status(r, http.StatusNoContent)
		return
	default:
		_ = render.Render(w, r, ErrServerError(err))
		return
	}
}

// DeleteUser passes id to storage to delete specific user
func (h handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.storage.DeleteUser(id)

	switch err {
	case i.UserNotFound:
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	case nil:
		render.Status(r, http.StatusNoContent)
		return
	default:
		_ = render.Render(w, r, ErrServerError(err))
		return
	}
}
