package v1

import (
	"github.com/go-chi/render"
	"net/http"
	"opencourse/database"
)

type RouteContext struct {
	DbContext database.DbContext
}

type OpenResponse[T any] struct {
	Payload *T     `json:"payload,omitempty"`
	Error   string `json:"error,omitempty"`
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

func WriteResponse[T any](writer http.ResponseWriter, request *http.Request, payload *T) {

	response := &OpenResponse[T]{Payload: payload}

	err := render.Render(writer, request, response)

	if err != nil {
		writer.WriteHeader(500)
	}
}

func WriteErrResponse[T any](writer http.ResponseWriter, request *http.Request, msg string, httpStatus int) {

	response := OpenResponse[T]{Error: msg}

	render.Status(request, httpStatus)
	err := render.Render(writer, request, &response)

	if err != nil {
		writer.WriteHeader(500)
	}
}
