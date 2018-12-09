package httpapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// getID extracts "id" parameter from path and converts it to uint.
func getID(req *http.Request) (uint, error) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// addID adds ID uint parameter to request context.
func addID(req *http.Request, id uint) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(int(id)))
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	return req.WithContext(ctx)
}
