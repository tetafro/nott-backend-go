package notepads

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

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/http/request"
	"github.com/tetafro/nott-backend-go/internal/testutils"
)

func TestController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	Uint := func(n uint) *uint {
		return &n
	}
	user := auth.User{ID: 1}

	t.Run("Get notepads", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notepads := []Notepad{
			{ID: 10, UserID: 20, FolderID: Uint(30), Title: "Notepad 10"},
			{ID: 15, UserID: 25, FolderID: Uint(35), Title: "Notepad 15"},
		}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Get(
			filter{userID: &user.ID},
		).Return(notepads, nil)

		c := NewController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = auth.AddUser(req, user)

		c.GetList(w, req)

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
			"data": [
				{
					"id": 10,
					"user_id": 20,
					"folder_id": 30,
					"title": "Notepad 10"
				},
				{
					"id": 15,
					"user_id": 25,
					"folder_id": 35,
					"title": "Notepad 15"
				}
			]
		}`)
	})

	t.Run("Failed to get notepads", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Get(
			filter{userID: &user.ID},
		).Return(nil, fmt.Errorf("error"))

		c := NewController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = auth.AddUser(req, user)

		c.GetList(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected status code %d, but got %d",
				http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("Create notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := Notepad{ID: id, UserID: user.ID, FolderID: Uint(30), Title: "Notepad 10"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Create(notepad).Return(notepad, nil)

		c := NewController(repoMock, log)

		payload, err := json.Marshal(notepad)
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}
		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		req = auth.AddUser(req, user)

		c.Create(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
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
				"id": 10,
				"user_id": 1,
				"folder_id": 30,
				"title": "Notepad 10"
			}
		}`)
	})

	t.Run("Failed to create notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := Notepad{ID: id, UserID: user.ID, FolderID: Uint(30), Title: "Notepad 10"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(notepad).Return(Notepad{}, fmt.Errorf("error"))

		c := NewController(repoMock, log)

		payload, err := json.Marshal(notepad)
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}
		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		req = auth.AddUser(req, user)
		req = request.AddUintPathParam(req, "id", id)

		c.Update(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected status code %d, but got %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("Get notepad by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepads := []Notepad{
			{ID: id, UserID: 20, FolderID: Uint(30), Title: "Notepad 10"},
		}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Get(
			filter{id: &id, userID: &user.ID},
		).Return(notepads, nil)

		c := NewController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = auth.AddUser(req, user)
		req = request.AddUintPathParam(req, "id", id)

		c.GetOne(w, req)

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
				"id": 10,
				"user_id": 20,
				"folder_id": 30,
				"title": "Notepad 10"
			}
		}`)
	})

	t.Run("Failed to get notepad by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Get(
			filter{id: &id, userID: &user.ID},
		).Return(nil, fmt.Errorf("error"))

		c := NewController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = auth.AddUser(req, user)
		req = request.AddUintPathParam(req, "id", id)

		c.GetOne(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected status code %d, but got %d",
				http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("Update notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := Notepad{ID: id, UserID: user.ID, FolderID: Uint(30), Title: "Notepad 10"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(notepad).Return(notepad, nil)

		c := NewController(repoMock, log)

		payload, err := json.Marshal(notepad)
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}
		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = auth.AddUser(req, user)
		req = request.AddUintPathParam(req, "id", id)

		c.Update(w, req)

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
				"id": 10,
				"user_id": 1,
				"folder_id": 30,
				"title": "Notepad 10"
			}
		}`)
	})

	t.Run("Failed to update notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := Notepad{ID: id, UserID: user.ID, FolderID: Uint(30), Title: "Notepad 10"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(notepad).Return(Notepad{}, fmt.Errorf("error"))

		c := NewController(repoMock, log)

		payload, err := json.Marshal(notepad)
		if err != nil {
			t.Fatalf("Failed to marshal json: %v", err)
		}
		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = auth.AddUser(req, user)
		req = request.AddUintPathParam(req, "id", id)

		c.Update(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected status code %d, but got %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("Delete notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := Notepad{ID: id, UserID: user.ID}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Delete(notepad).Return(nil)

		c := NewController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req = auth.AddUser(req, user)
		req = request.AddUintPathParam(req, "id", id)

		c.Delete(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status code %d, but got %d", http.StatusNoContent, resp.StatusCode)
		}
	})

	t.Run("Failed to delete notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := Notepad{ID: id, UserID: user.ID}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Delete(notepad).Return(fmt.Errorf("error"))

		c := NewController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, nil)
		req = auth.AddUser(req, user)
		req = request.AddUintPathParam(req, "id", id)

		c.Delete(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected status code %d, but got %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})
}
