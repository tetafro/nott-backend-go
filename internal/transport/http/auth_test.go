package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

func TestAuthController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	password := "qwerty"
	hash := "$2a$14$EsnwEn3C6cxQUWXvUpJ6S.XsJSku11hTSULXn8NEIG1diGcGEgrii"
	user := auth.User{ID: 1, Email: "bob@example.com", Password: hash}
	token := auth.Token{ID: 10, UserID: user.ID, String: "qwerty123", TTL: 10}

	t.Run("Succesful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(user, nil)

		tokensRepoMock := storage.NewMockTokensRepo(ctrl)
		tokensRepoMock.EXPECT().Create(
			auth.Token{UserID: token.UserID},
		).Return(token, nil)

		c := NewAuthController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: password})
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Login(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
			"data": {
				"string": "qwerty123",
				"ttl": 10
			}
		}`)
	})

	t.Run("Failed login because of users repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(auth.User{}, fmt.Errorf("error"))

		tokensRepoMock := storage.NewMockTokensRepo(ctrl)

		c := NewAuthController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: password})
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Login(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Failed login because of tokens repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(user, nil)

		tokensRepoMock := storage.NewMockTokensRepo(ctrl)
		tokensRepoMock.EXPECT().Create(
			auth.Token{UserID: token.UserID},
		).Return(auth.Token{}, fmt.Errorf("error"))

		c := NewAuthController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: password})
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Login(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	// "Get profile" is not tested here, because the only method,
	// that can be tested, works inside auth middleware

	t.Run("Update profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := auth.User{ID: 10, Email: "bob@example.com"}

		userRepoMock := storage.NewMockUsersRepo(ctrl)
		userRepoMock.EXPECT().Update(user).Return(user, nil)

		tokenRepoMock := storage.NewMockTokensRepo(ctrl)

		c := NewAuthController(userRepoMock, tokenRepoMock, log)

		payload, err := json.Marshal(user)
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = addUser(req, auth.User{ID: 10, Email: "alice@example.com"})

		c.UpdateProfile(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
			"data": {
				"id": 10,
				"email": "bob@example.com"
			}
		}`)
	})

	t.Run("Failed to update profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := auth.User{Email: "bob@example.com"}

		userRepoMock := storage.NewMockUsersRepo(ctrl)
		userRepoMock.EXPECT().Update(user).Return(auth.User{}, fmt.Errorf("fatal error"))

		tokenRepoMock := storage.NewMockTokensRepo(ctrl)

		c := NewAuthController(userRepoMock, tokenRepoMock, log)

		payload, err := json.Marshal(user)
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = addUser(req, auth.User{Email: "alice@example.com"})

		c.UpdateProfile(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}