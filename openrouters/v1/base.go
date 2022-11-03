package v1

import (
	"net/http"
	"opencourse/database"
)

type RouteContext struct {
	DbContext database.DbContext
}

type OpenResponse[T any] struct {
	Data T `json:"data"`
}

func (or *OpenResponse[T]) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
