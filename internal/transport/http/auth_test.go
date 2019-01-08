package httpapi

import (
	"bytes"
	"encoding/json"
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

func TestAuthController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	password := "qwerty"
	hash := "$2a$14$EsnwEn3C6cxQUWXvUpJ6S.XsJSku11hTSULXn8NEIG1diGcGEgrii"
	user := auth.User{ID: 1, Email: "bob@example.com", Password: hash}
	token := auth.Token{AccessToken: "qwerty123", ExpiresAt: 10}

	t.Run("Succesful registration", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(auth.User{}, domain.ErrNotFound)
		usersRepoMock.EXPECT().Create(gomock.Any()).Return(user, nil)

		tokenerMock := auth.NewMockTokener(ctrl)
		tokenerMock.EXPECT().Issue(user).Return(token, nil)

		c := NewAuthController(usersRepoMock, tokenerMock, log)

		payload, err := json.Marshal(authRequest{Email: user.Email, Password: password})
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Register(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
			"data": {
				"access_token": "qwerty123",
				"expires_at": 10
			}
		}`)
	})

	t.Run("Failed registration because user already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(auth.User{ID: 1}, nil)

		tokenerMock := auth.NewMockTokener(ctrl)

		c := NewAuthController(usersRepoMock, tokenerMock, log)

		payload, err := json.Marshal(authRequest{Email: user.Email, Password: password})
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Register(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
	})

	t.Run("Failed registration because of users repo error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(auth.User{}, domain.ErrNotFound)
		usersRepoMock.EXPECT().Create(gomock.Any()).Return(auth.User{}, errors.New("error"))

		tokenerMock := auth.NewMockTokener(ctrl)

		c := NewAuthController(usersRepoMock, tokenerMock, log)

		payload, err := json.Marshal(authRequest{Email: user.Email, Password: password})
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Register(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Succesful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(user, nil)

		tokenerMock := auth.NewMockTokener(ctrl)
		tokenerMock.EXPECT().Issue(user).Return(token, nil)

		c := NewAuthController(usersRepoMock, tokenerMock, log)

		payload, err := json.Marshal(authRequest{Email: user.Email, Password: password})
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
				"access_token": "qwerty123",
				"expires_at": 10
			}
		}`)
	})

	t.Run("Failed login because of users repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(auth.User{}, errors.New("error"))

		tokenerMock := auth.NewMockTokener(ctrl)

		c := NewAuthController(usersRepoMock, tokenerMock, log)

		payload, err := json.Marshal(authRequest{Email: user.Email, Password: password})
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Login(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Get profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := auth.User{ID: 10, Email: "bob@example.com"}

		userRepoMock := storage.NewMockUsersRepo(ctrl)
		userRepoMock.EXPECT().GetByID(user.ID).Return(user, nil)

		tokenerMock := auth.NewMockTokener(ctrl)

		c := NewAuthController(userRepoMock, tokenerMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUserID(req, user.ID)

		c.GetProfile(w, req)

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

	t.Run("Update profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := auth.User{ID: 10, Email: "bob@example.com"}

		userRepoMock := storage.NewMockUsersRepo(ctrl)
		userRepoMock.EXPECT().Update(user).Return(user, nil)

		tokenerMock := auth.NewMockTokener(ctrl)

		c := NewAuthController(userRepoMock, tokenerMock, log)

		payload, err := json.Marshal(user)
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = addUserID(req, user.ID)

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
		userRepoMock.EXPECT().Update(user).Return(auth.User{}, errors.New("error"))

		tokenerMock := auth.NewMockTokener(ctrl)

		c := NewAuthController(userRepoMock, tokenerMock, log)

		payload, err := json.Marshal(user)
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = addUserID(req, user.ID)

		c.UpdateProfile(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}
