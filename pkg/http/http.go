package http

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
)

type Handler func(*http.Request) Response

type Response interface {
	StatusCode() int
}

type Success struct {
	Status     string      `json:"status"`
	Payload    interface{} `json:"payload,omitempty"`
	statusCode int
}

func (s Success) StatusCode() int {
	return s.statusCode
}

func OK(payload interface{}) *Success {
	return &Success{
		Status:     StatusOk,
		statusCode: http.StatusOK,
		Payload:    payload,
	}
}

type Error struct {
	Status     string      `json:"status"`
	Payload    interface{} `json:"payload,omitempty"`
	err        error
	statusCode int
}

func (e Error) StatusCode() int {
	return e.statusCode
}

type ErrorPayload struct {
	Message string `json:"message"`
}

func BadRequest(err error) *Error {
	return &Error{
		Status:     StatusError,
		statusCode: http.StatusBadRequest,
		Payload: ErrorPayload{
			Message: err.Error(),
		},
		err: err,
	}
}

func NotFound(err error) *Error {
	return &Error{
		Status:     StatusError,
		statusCode: http.StatusNotFound,
		Payload: ErrorPayload{
			Message: err.Error(),
		},
		err: err,
	}
}

func AddHandler(
	mountMethod func(pattern string, h http.HandlerFunc),
	pattern string,
	handler Handler,
) {
	mountMethod(pattern, converter(handler))
}

func converter(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := h(r)
		if resp == nil {
			return
		}

		writeResponse(w, resp)
	}
}

func writeResponse(w http.ResponseWriter, response Response) {
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(response.StatusCode())

	_, _ = w.Write(b)
}
