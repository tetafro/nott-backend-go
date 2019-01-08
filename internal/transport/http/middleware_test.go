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

func TestAuthMiddleware(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	user := auth.User{ID: 10, Email: "bob@example.com"}

	t.Run("Authorize user by token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByToken("token").Return(user, nil)

		mw := NewAuthMiddleware(usersRepoMock, log)

		h := func(w http.ResponseWriter, r *http.Request) {
			// Check user in request context
			u := getUser(r)
			assert.Equal(t, user, *u)

			w.Write([]byte("ok")) // nolint
		}
		ts := httptest.NewServer(mw(http.HandlerFunc(h)))
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		assert.NoError(t, err)

		req.Header.Add("Authorization", `Token token="token"`)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Fail to authorize user with no token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)

		mw := NewAuthMiddleware(usersRepoMock, log)

		h := func(w http.ResponseWriter, r *http.Request) {
			// Won't get here
			w.Write([]byte("ok")) // nolint
		}
		ts := httptest.NewServer(mw(http.HandlerFunc(h)))
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Fail to authorize user with invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByToken("wrong-token").Return(auth.User{}, domain.ErrNotFound)

		mw := NewAuthMiddleware(usersRepoMock, log)

		h := func(w http.ResponseWriter, r *http.Request) {
			// Won't get here
			w.Write([]byte("ok")) // nolint
		}
		ts := httptest.NewServer(mw(http.HandlerFunc(h)))
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		assert.NoError(t, err)

		req.Header.Add("Authorization", `Token token="wrong-token"`)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Fail to authorize user with malformed token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)

		mw := NewAuthMiddleware(usersRepoMock, log)

		h := func(w http.ResponseWriter, r *http.Request) {
			// Won't get here
			w.Write([]byte("ok")) // nolint
		}
		ts := httptest.NewServer(mw(http.HandlerFunc(h)))
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		assert.NoError(t, err)

		req.Header.Add("Authorization", `Token "token"`)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Fail to authorize user because of users repo error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := storage.NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByToken("token").Return(auth.User{}, errors.New("error"))

		mw := NewAuthMiddleware(usersRepoMock, log)

		h := func(w http.ResponseWriter, r *http.Request) {
			// Won't get here
			w.Write([]byte("ok")) // nolint
		}
		ts := httptest.NewServer(mw(http.HandlerFunc(h)))
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		assert.NoError(t, err)

		req.Header.Add("Authorization", `Token token="token"`)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
