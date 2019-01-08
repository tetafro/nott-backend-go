package httpapi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

func TestOAuthController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	t.Run("Get list of providers", func(t *testing.T) {
		c := NewOAuthController(
			map[string]*auth.OAuthProvider{
				"example": {Name: "example", URL: "http://example.com"},
			},
			nil, nil, nil,
		)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)

		c.Providers(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
			"data": [
				{
					"name": "example",
					"url": "http://example.com"
				}
			]
		}`)
	})

	t.Run("Handle existing user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := auth.User{ID: 10, Email: "bob@example.com"}

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(user, nil)

		tokenerMock := auth.NewMockTokener(ctrl)
		tokenerMock.EXPECT().Issue(user).Return(
			auth.Token{AccessToken: "qwerty", ExpiresAt: 10}, nil,
		)

		c := NewOAuthController(nil, usersRepoMock, tokenerMock, log)

		token, err := c.handleUser(user.Email)
		assert.NoError(t, err)
		assert.Equal(t, "qwerty", token.AccessToken)
	})

	t.Run("Handle new user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := auth.User{ID: 10, Email: "bob@example.com"}

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(auth.User{}, domain.ErrNotFound)
		usersRepoMock.EXPECT().Create(auth.User{Email: user.Email}).Return(user, nil)

		tokenerMock := auth.NewMockTokener(ctrl)
		tokenerMock.EXPECT().Issue(user).Return(
			auth.Token{AccessToken: "qwerty", ExpiresAt: 10}, nil,
		)

		c := NewOAuthController(nil, usersRepoMock, tokenerMock, log)

		token, err := c.handleUser(user.Email)
		assert.NoError(t, err)
		assert.Equal(t, "qwerty", token.AccessToken)
	})

	t.Run("Fail to handle user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := auth.User{ID: 10, Email: "bob@example.com"}

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(auth.User{}, errors.New("error"))

		tokenerMock := auth.NewMockTokener(ctrl)

		c := NewOAuthController(nil, usersRepoMock, tokenerMock, log)

		_, err := c.handleUser(user.Email)
		assert.Error(t, err)
	})
}
