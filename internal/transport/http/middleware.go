package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// userKey is a key for user structure inside request context.
type userKey struct{}

// NewAuthMiddleware creates middleware that authenticates users.
func NewAuthMiddleware(users storage.UsersRepo, log logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := getToken(req)
			if token == "" {
				unauthorized(w)
				return
			}
			u, err := users.GetByToken(token)
			if err == domain.ErrNotFound {
				unauthorized(w)
				return
			}
			if err != nil {
				log.Errorf("Failed to get user from database: %v", err)
				internalServerError(w)
				return
			}
			req = addUser(req, u)
			next.ServeHTTP(w, req)
		})
	}
}

// addUser adds user to request context.
func addUser(req *http.Request, u auth.User) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, userKey{}, u)
	return req.WithContext(ctx)
}

// getUser extracts user from request context.
func getUser(req *http.Request) *auth.User {
	val := req.Context().Value(userKey{})
	if val == nil {
		panic("Failed to get user from request context")
	}
	u, ok := val.(auth.User)
	if !ok {
		panic("Unknown user value in request context")
	}
	return &u
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
