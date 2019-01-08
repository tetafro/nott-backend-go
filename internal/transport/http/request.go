package httpapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// getID extracts "id" parameter from path.
func getID(req *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(req, "id"))
}

// addID adds ID int parameter to request context.
func addID(req *http.Request, id int) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(id))
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	return req.WithContext(ctx)
}
