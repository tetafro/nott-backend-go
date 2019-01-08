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
)

func TestAddUser(t *testing.T) {
	userID := 10

	req := &http.Request{}
	req = addUserID(req, userID)

	reqUserID := req.Context().Value(userIDKey{})
	assert.Equal(t, userID, reqUserID)
}

func TestGetUser(t *testing.T) {
	userID := 10

	req := &http.Request{}
	req = addUserID(req, userID)

	reqUserID := getUserID(req)
	assert.Equal(t, userID, reqUserID)
}

func TestAuthMiddleware(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	user := auth.User{ID: 10, Email: "bob@example.com"}

	t.Run("Authorize user by token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenerMock := auth.NewMockTokener(ctrl)
		tokenerMock.EXPECT().Parse("qwerty").Return(user.ID, nil)

		mw := NewAuthMiddleware(tokenerMock, log)

		h := func(w http.ResponseWriter, r *http.Request) {
			// Check user in request context
			assert.Equal(t, user.ID, getUserID(r))

			w.Write([]byte("ok")) // nolint
		}
		ts := httptest.NewServer(mw(http.HandlerFunc(h)))
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
		assert.NoError(t, err)

		req.Header.Add("Authorization", `Token token="qwerty"`)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Fail to authorize user with no token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenerMock := auth.NewMockTokener(ctrl)

		mw := NewAuthMiddleware(tokenerMock, log)

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

	t.Run("Fail to authorize user with malformed token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenerMock := auth.NewMockTokener(ctrl)

		mw := NewAuthMiddleware(tokenerMock, log)

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

	t.Run("Fail to authorize user with invalid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		tokenerMock := auth.NewMockTokener(ctrl)
		tokenerMock.EXPECT().Parse("wrong-token").Return(0, errors.New("error"))

		mw := NewAuthMiddleware(tokenerMock, log)

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
}
