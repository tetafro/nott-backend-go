package httpapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetafro/nott-backend-go/internal/auth"
)

func TestAddUser(t *testing.T) {
	user := auth.User{ID: 10}

	req := &http.Request{}
	req = addUser(req, user)

	reqUser := req.Context().Value(userKey{})
	assert.Equal(t, user, reqUser)
}

func TestGetUser(t *testing.T) {
	user := auth.User{ID: 10}

	req := &http.Request{}
	req = addUser(req, user)

	reqUser := getUser(req)
	assert.Equal(t, user, *reqUser)
}
