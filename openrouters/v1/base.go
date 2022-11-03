package v1

import (
	"net/http"
	"opencourse/database"
)

type RouteContext struct {
	DbContext database.DbContext
}

type OpenResponse[T any] struct {
	Payload T `json:"payload"`
}

type OpenRequest[T any] struct {
	Payload T `json:"payload"`
}

func (or *OpenResponse[T]) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (or *OpenRequest[T]) Bind(r *http.Request) error {
	return nil
}
