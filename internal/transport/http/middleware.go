package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/auth"
)

// userIDKey is a key for user id value inside request context.
type userIDKey struct{}

// NewAuthMiddleware creates middleware that authenticates users.
func NewAuthMiddleware(tokener auth.Tokener, log logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := getToken(req)
			if token == "" {
				unauthorized(w)
				return
			}
			id, err := tokener.Parse(token)
			if err != nil {
				unauthorized(w)
				return
			}
			req = addUserID(req, id)
			next.ServeHTTP(w, req)
		})
	}
}

// addUserID adds user id to request context.
func addUserID(req *http.Request, id int) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, userIDKey{}, id)
	return req.WithContext(ctx)
}

// getUserID extracts user id from request context.
func getUserID(req *http.Request) int {
	val := req.Context().Value(userIDKey{})
	if val == nil {
		panic("Failed to get user id from request context")
	}
	id, ok := val.(int)
	if !ok {
		panic("Invalid user id type in request context")
	}
	return id
}

// getToken gets token from HTTP header.
// Header format (RFC2617):
// Authorization: Token token="abcd1234"
func getToken(req *http.Request) string {
	auth, ok := req.Header["Authorization"]
	if !ok || len(auth) == 0 {
		return ""
	}

	token := auth[0]
	if !strings.HasPrefix(token, `Token token="`) || !strings.HasSuffix(token, `"`) {
		return ""
	}

	return token[13 : len(token)-1]
}
