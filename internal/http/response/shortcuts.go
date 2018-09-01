package response

import (
	"fmt"
	"net/http"
)

// NotFound is a shortcut for 404 response with a
// default message.
func NotFound() Response {
	code := http.StatusNotFound
	err := fmt.Errorf(http.StatusText(code))
	return New().WithStatus(code).WithError(err)
}

// Unauthorized is a shortcut for 401 response with a
// default message.
func Unauthorized() Response {
	code := http.StatusUnauthorized
	err := fmt.Errorf(http.StatusText(code))
	return New().WithStatus(code).WithError(err)
}

// InternalServerError is a shortcut for 500 response with a
// default message.
func InternalServerError() Response {
	code := http.StatusInternalServerError
	err := fmt.Errorf(http.StatusText(code))
	return New().WithStatus(code).WithError(err)
}
