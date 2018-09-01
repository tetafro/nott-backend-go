package auth

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/http/request"
	"github.com/tetafro/nott-backend-go/internal/http/response"
)

// UserKey is a key for user structure inside request context.
type UserKey struct{}

// NewAuthMiddleware creates middleware that authenticates users.
func NewAuthMiddleware(users UsersRepo, log logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := request.GetToken(req)
			if token == "" {
				response.Unauthorized().Write(w)
				return
			}
			u, err := users.GetByToken(token)
			if err == ErrNotFound {
				response.Unauthorized().Write(w)
				return
			}
			if err != nil {
				log.Errorf("Failed to get user from database: %v", err)
				return
			}
			req = AddUser(req, u)
			next.ServeHTTP(w, req)
		})
	}
}

// AddUser adds user to request context.
func AddUser(req *http.Request, u User) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, UserKey{}, u)
	return req.WithContext(ctx)
}

// GetUser extracts user from request context.
func GetUser(req *http.Request) *User {
	val := req.Context().Value(UserKey{})
	if val == nil {
		panic("Failed to get user from request context")
	}
	u, ok := val.(User)
	if !ok {
		panic("Unknown user value in request context")
	}
	return &u
}
