package httpapi

import (
	"encoding/json"
	"net/http"
)

const contentType = "application/json; charset=utf-8"

type dataResponse struct {
	Data interface{} `json:"data"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func respond(w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)

	var resp interface{}
	if code >= 200 && code <= 299 {
		resp = dataResponse{Data: data}
	} else {
		resp = errorResponse{Error: data.(string)}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func badRequest(w http.ResponseWriter, err string) {
	respond(w, http.StatusBadRequest, err)
}

func notFound(w http.ResponseWriter) {
	respond(w, http.StatusNotFound, "not found")
}

func unauthorized(w http.ResponseWriter) {
	respond(w, http.StatusUnauthorized, "unauthorized")
}

func internalServerError(w http.ResponseWriter) {
	respond(w, http.StatusInternalServerError, "internal server error")
}
