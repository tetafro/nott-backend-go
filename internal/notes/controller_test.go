package notes

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
	"github.com/tetafro/nott-backend-go/internal/httpx/request"
	"github.com/tetafro/nott-backend-go/internal/testutils"
)

func TestController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	Uint := func(n uint) *uint {
		return &n
	}
	user := auth.User{ID: 1}

	t.Run("Get notes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notes := []Note{
			{ID: 10, UserID: user.ID, NotepadID: Uint(30), Title: "Note 10", Text: "Hello"},
			{ID: 15, UserID: user.ID, NotepadID: Uint(30), Title: "Note 15", Text: "Hello"},
		}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Get(
			filter{userID: &user.ID, notepadID: Uint(10)},
		).Return(notes, nil)

		c := NewController(repoMock, log)

		url := "/?notepad_id=10"
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
					"user_id": 1,
					"notepad_id": 30,
					"title": "Note 10",
					"text": "Hello"
				},
				{
					"id": 15,
					"user_id": 1,
					"notepad_id": 30,
					"title": "Note 15",
					"text": "Hello"
				}
			]
		}`)
	})

	t.Run("Failed to get notes", func(t *testing.T) {
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

	t.Run("Create note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		note := Note{ID: id, UserID: user.ID, NotepadID: Uint(30), Title: "Note 10", Text: "Hello"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Create(note).Return(note, nil)

		c := NewController(repoMock, log)

		payload, err := json.Marshal(note)
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
				"notepad_id": 30,
				"title": "Note 10",
				"text": "Hello"
			}
		}`)
	})

	t.Run("Failed to create note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		note := Note{ID: id, UserID: user.ID, NotepadID: Uint(30), Title: "Note 10", Text: "Hello"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(note).Return(Note{}, fmt.Errorf("error"))

		c := NewController(repoMock, log)

		payload, err := json.Marshal(note)
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

	t.Run("Get note by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notes := []Note{
			{ID: id, UserID: user.ID, NotepadID: Uint(30), Title: "Note 10", Text: "Hello"},
		}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Get(
			filter{id: &id, userID: &user.ID},
		).Return(notes, nil)

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
				"user_id": 1,
				"notepad_id": 30,
				"title": "Note 10",
				"text": "Hello",
				"html": "\u003cp\u003eHello\u003c/p\u003e\n"
			}
		}`)
	})

	t.Run("Failed to get note by id", func(t *testing.T) {
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

	t.Run("Update note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		note := Note{ID: id, UserID: user.ID, NotepadID: Uint(30), Title: "Note 10", Text: "Hello"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(note).Return(note, nil)

		c := NewController(repoMock, log)

		payload, err := json.Marshal(note)
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
				"notepad_id": 30,
				"title": "Note 10",
				"text": "Hello"
			}
		}`)
	})

	t.Run("Failed to update note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		note := Note{ID: id, UserID: user.ID, NotepadID: Uint(30), Title: "Note 10", Text: "Hello"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(note).Return(Note{}, fmt.Errorf("error"))

		c := NewController(repoMock, log)

		payload, err := json.Marshal(note)
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

	t.Run("Delete note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		note := Note{ID: id, UserID: user.ID}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Delete(note).Return(nil)

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

	t.Run("Failed to delete note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		note := Note{ID: id, UserID: user.ID}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Delete(note).Return(fmt.Errorf("error"))

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
