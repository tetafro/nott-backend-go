package auth

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

	"github.com/tetafro/nott-backend-go/internal/testutils"
)

func TestController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	user := User{ID: 1, Email: "bob@example.com", Password: "qwerty"}
	token := Token{ID: 10, UserID: user.ID, String: "qwerty123", TTL: 10}

	t.Run("Succesful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmailAndPassword(
			user.Email, user.Password,
		).Return(user, nil)

		tokensRepoMock := NewMockTokensRepo(ctrl)
		tokensRepoMock.EXPECT().Create(
			Token{UserID: token.UserID},
		).Return(token, nil)

		c := NewController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: user.Password})
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Login(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		err = resp.Body.Close()
		if err != nil {
			t.Fatalf("Failed to close response body: %v", err)
		}
		testutils.AssertResponse(t, string(body), `{
			"data": {
				"string": "qwerty123",
				"ttl": 10,
				"created": "0001-01-01T00:00:00Z"
			}
		}`)
	})

	t.Run("Failed login because of users repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmailAndPassword(
			user.Email, user.Password,
		).Return(User{}, fmt.Errorf("error"))

		tokensRepoMock := NewMockTokensRepo(ctrl)

		c := NewController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: user.Password})
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Login(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected status code %d, but got %d",
				http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("Failed login because of tokens repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmailAndPassword(
			user.Email, user.Password,
		).Return(user, nil)

		tokensRepoMock := NewMockTokensRepo(ctrl)
		tokensRepoMock.EXPECT().Create(
			Token{UserID: token.UserID},
		).Return(Token{}, fmt.Errorf("error"))

		c := NewController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: user.Password})
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

		c.Login(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected status code %d, but got %d",
				http.StatusInternalServerError, resp.StatusCode)
		}
	})
}
