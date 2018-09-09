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

	password := "qwerty"
	hash := "$2a$14$EsnwEn3C6cxQUWXvUpJ6S.XsJSku11hTSULXn8NEIG1diGcGEgrii"
	user := User{ID: 1, Email: "bob@example.com", Password: hash}
	token := Token{ID: 10, UserID: user.ID, String: "qwerty123", TTL: 10}

	t.Run("Succesful login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(user, nil)

		tokensRepoMock := NewMockTokensRepo(ctrl)
		tokensRepoMock.EXPECT().Create(
			Token{UserID: token.UserID},
		).Return(token, nil)

		c := NewController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: password})
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
		testutils.AssertJSON(t, string(body), `{
			"data": {
				"string": "qwerty123",
				"ttl": 10
			}
		}`)
	})

	t.Run("Failed login because of users repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepoMock := NewMockUsersRepo(ctrl)
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(User{}, fmt.Errorf("error"))

		tokensRepoMock := NewMockTokensRepo(ctrl)

		c := NewController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: password})
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
		usersRepoMock.EXPECT().GetByEmail(user.Email).Return(user, nil)

		tokensRepoMock := NewMockTokensRepo(ctrl)
		tokensRepoMock.EXPECT().Create(
			Token{UserID: token.UserID},
		).Return(Token{}, fmt.Errorf("error"))

		c := NewController(usersRepoMock, tokensRepoMock, log)

		payload, err := json.Marshal(loginRequest{Email: user.Email, Password: password})
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

	// "Get profile" is not tested here, because the only method,
	// that can be tested, works inside auth middleware

	t.Run("Update profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := User{ID: 10, Email: "bob@example.com"}

		userRepoMock := NewMockUsersRepo(ctrl)
		userRepoMock.EXPECT().Update(user).Return(user, nil)

		tokenRepoMock := NewMockTokensRepo(ctrl)

		c := NewController(userRepoMock, tokenRepoMock, log)

		payload, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}
		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = AddUser(req, User{ID: 10, Email: "alice@example.com"})

		c.UpdateProfile(w, req)

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
		testutils.AssertJSON(t, string(body), `{
			"data": {
				"id": 10,
				"email": "bob@example.com"
			}
		}`)
	})

	t.Run("Failed to update profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := User{Email: "bob@example.com"}

		userRepoMock := NewMockUsersRepo(ctrl)
		userRepoMock.EXPECT().Update(user).Return(User{}, fmt.Errorf("fatal error"))

		tokenRepoMock := NewMockTokensRepo(ctrl)

		c := NewController(userRepoMock, tokenRepoMock, log)

		payload, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}
		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = AddUser(req, User{Email: "alice@example.com"})

		c.UpdateProfile(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected status code %d, but got %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})
}
