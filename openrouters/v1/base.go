package v1

import (
	"encoding/json"
	"github.com/go-chi/httplog"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"net/http"
	"opencourse/database"
)

const (
	ErrInternal          = 1 // ErrInternal - internal business logic
	ErrBinding           = 2 // ErrBinding - data binding error
	ErrParameter         = 3 // ErrParameter - parameter wrong error
	ErrLoginOrPassword   = 4 // ErrLoginOrPassword - login or password is incorrect
	ErrUserAlreadyExists = 5 // ErrUserAlreadyExists - user with same login already exist
	ErrValid             = 6 // ErrValid - validation error
)

// ResponseError model with error description
type ResponseError struct {
	Message string `json:"message"` // Message error
	Code    int    `json:"code"`    // Code error
}

// RouteContext contains data for request handlers
type RouteContext struct {
	DbContext database.DbContext // DbContext, contains methods and properties for work with db
	TokenAuth *jwtauth.JWTAuth   // TokenAuth contains methods for decode and encode jwt tokens
}

// Response is model for http handler response. Contains properties with user data and error
type Response[T any] struct {
	Payload *T             `json:"payload,omitempty"` // Payload is a user model with data
	Error   *ResponseError `json:"error,omitempty"`   // Error field contains a description of the error
}

// Request is a user request model
type Request[T any] struct {
	Payload T `json:"payload"` //Payload is a user data in request
}

// Render method for render response. It's empty.
func (or *Response[T]) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind method for data binding in request. It's empty.
func (or *Request[T]) Bind(r *http.Request) error {
	return nil
}

// WriteResponse function for build and write data for success response
func WriteResponse[T any](writer http.ResponseWriter, request *http.Request, payload *T) {

	response := &Response[T]{Payload: payload}

	err := render.Render(writer, request, response)

	if err != nil {
		writer.WriteHeader(500)
	}
}

// WriteErrResponse function for build and write data for error response
func WriteErrResponse(writer http.ResponseWriter, request *http.Request, internalError error, responseError *ResponseError, httpStatus int) {

	if internalError != nil {
		httplog.LogEntrySetField(request.Context(), "internal_error", internalError.Error())
	}

	if responseError != nil {
		j, err := json.Marshal(responseError)

		if err != nil {
			writer.WriteHeader(500)
			httplog.LogEntrySetField(request.Context(), "internal_error", err.Error())
			return
		}

		httplog.LogEntrySetField(request.Context(), "response_error", string(j))
	}

	response := Response[string]{Error: responseError}

	render.Status(request, httpStatus)
	err := render.Render(writer, request, &response)

	if err != nil {
		writer.WriteHeader(500)
		httplog.LogEntrySetField(request.Context(), "err.render", err.Error())
	}
}
