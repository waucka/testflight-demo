package api

import (
	"net/http"
)

type ErrorDescription struct {
	Status int
	Error string
}

func InternalError(msg string) {
	panic(&ErrorDescription{
		Status: http.StatusInternalServerError,
		Error: msg,
	})
}

func Unauthorized(msg string) {
	panic(&ErrorDescription{
		Status: http.StatusUnauthorized,
		Error: msg,
	})
}

func Forbidden(msg string) {
	panic(&ErrorDescription{
		Status: http.StatusForbidden,
		Error: msg,
	})
}

func NotFound(msg string) {
	panic(&ErrorDescription{
		Status: http.StatusNotFound,
		Error: msg,
	})
}

func BadRequest(msg string) {
	panic(&ErrorDescription{
		Status: http.StatusBadRequest,
		Error: msg,
	})
}

func Gone(msg string) {
	panic(&ErrorDescription{
		Status: http.StatusGone,
		Error: msg,
	})
}
