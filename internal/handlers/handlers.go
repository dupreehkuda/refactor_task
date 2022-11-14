package handlers

import i "github.com/dupreehkuda/refactor_task/internal/interfaces"

type handlers struct {
	storage i.Storer
}

// New creates new instance of handlers
func New(storage i.Storer) *handlers {
	return &handlers{storage: storage}
}
