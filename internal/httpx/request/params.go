package request

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// GetUintPathParam extracts named parameter from path
// and converts it to uint.
func GetUintPathParam(req *http.Request, key string) (uint, error) {
	i, err := strconv.Atoi(chi.URLParam(req, key))
	if err != nil {
		return 1, err
	}
	return uint(i), nil
}

// AddUintPathParam adds uint parameter to request context.
func AddUintPathParam(req *http.Request, key string, value uint) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, strconv.Itoa(int(value)))
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	return req.WithContext(ctx)
}
