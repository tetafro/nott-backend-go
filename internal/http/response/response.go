package response

import (
	"encoding/json"
	"net/http"
)

const contentType = "application/json; charset=utf-8"

// Response is an HTTP JSON response.
type Response struct {
	// StatusCode is the response HTTP status code.
	StatusCode int

	// Data is the document's “primary data”
	Data interface{}

	// Error is an error message
	Error error
}

// New returns a new Response.
// By default it has a status code 200.
func New() Response {
	return Response{StatusCode: http.StatusOK}
}

// WithStatus returns response with status code added.
func (r Response) WithStatus(code int) Response {
	r.StatusCode = code
	return r
}

// WithData returns response with data added.
func (r Response) WithData(data interface{}) Response {
	r.Data = data
	return r
}

// WithError returns response with error added.
func (r Response) WithError(err error) Response {
	r.Error = err
	return r
}

// Write dispatches the response.
func (r Response) Write(w http.ResponseWriter) {
	if r.StatusCode == http.StatusNoContent {
		w.WriteHeader(r.StatusCode)
		return
	}

	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(r.StatusCode)

	enc := json.NewEncoder(w)

	var err error
	switch {
	case r.StatusCode >= 200 && r.StatusCode <= 299:
		err = enc.Encode(struct {
			Data interface{} `json:"data"`
		}{
			Data: r.Data,
		})
	default:
		err = enc.Encode(struct {
			Error string `json:"error"`
		}{
			Error: r.Error.Error(),
		})
	}

	if err != nil {
		panic(err)
	}
}
