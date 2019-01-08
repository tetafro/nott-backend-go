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
	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

func TestNotesNController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	Int := func(n int) *int {
		return &n
	}
	user := auth.User{ID: 1}

	t.Run("Get notes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notes := []domain.Note{
			{ID: 10, UserID: user.ID, NotepadID: 30, Title: "Note 10", Text: "Hello"},
			{ID: 15, UserID: user.ID, NotepadID: 30, Title: "Note 15", Text: "Hello"},
		}

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotesFilter{UserID: &user.ID, NotepadID: Int(10)},
		).Return(notes, nil)

		c := NewNotesController(repoMock, log)

		url := "/?notepad_id=10"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)

		c.GetList(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
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

	t.Run("Fail to get notes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotesFilter{UserID: &user.ID},
		).Return(nil, fmt.Errorf("error"))

		c := NewNotesController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)

		c.GetList(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Create note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		note := domain.Note{ID: id, UserID: user.ID, NotepadID: 30, Title: "Note 10", Text: "Hello"}

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Create(note).Return(note, nil)

		c := NewNotesController(repoMock, log)

		payload, err := json.Marshal(note)
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		req = addUser(req, user)

		c.Create(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusCreated)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
			"data": {
				"id": 10,
				"user_id": 1,
				"notepad_id": 30,
				"title": "Note 10",
				"text": "Hello"
			}
		}`)
	})

	t.Run("Fail to create note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		note := domain.Note{ID: id, UserID: user.ID, NotepadID: 30, Title: "Note 10", Text: "Hello"}

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Update(note).Return(domain.Note{}, fmt.Errorf("error"))

		c := NewNotesController(repoMock, log)

		payload, err := json.Marshal(note)
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		req = addUser(req, user)
		req = addID(req, id)

		c.Update(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Get note by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		notes := []domain.Note{
			{ID: id, UserID: user.ID, NotepadID: 30, Title: "Note 10", Text: "Hello"},
		}

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotesFilter{ID: &id, UserID: &user.ID},
		).Return(notes, nil)

		c := NewNotesController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.GetOne(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
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

	t.Run("Fail to get note by non-existing id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotesFilter{ID: &id, UserID: &user.ID},
		).Return(nil, nil)

		c := NewNotesController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.GetOne(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	})

	t.Run("Fail to get note by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotesFilter{ID: &id, UserID: &user.ID},
		).Return(nil, fmt.Errorf("error"))

		c := NewNotesController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.GetOne(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Update note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		note := domain.Note{ID: id, UserID: user.ID, NotepadID: 30, Title: "Note 10", Text: "Hello"}

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Update(note).Return(note, nil)

		c := NewNotesController(repoMock, log)

		payload, err := json.Marshal(note)
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = addUser(req, user)
		req = addID(req, id)

		c.Update(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = resp.Body.Close()
		assert.NoError(t, err)

		assert.JSONEq(t, string(body), `{
			"data": {
				"id": 10,
				"user_id": 1,
				"notepad_id": 30,
				"title": "Note 10",
				"text": "Hello"
			}
		}`)
	})

	t.Run("Fail to update note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		note := domain.Note{ID: id, UserID: user.ID, NotepadID: 30, Title: "Note 10", Text: "Hello"}

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Update(note).Return(domain.Note{}, fmt.Errorf("error"))

		c := NewNotesController(repoMock, log)

		payload, err := json.Marshal(note)
		assert.NoError(t, err)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req = addUser(req, user)
		req = addID(req, id)

		c.Update(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Delete note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		note := domain.Note{ID: id, UserID: user.ID}

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Delete(note).Return(nil)

		c := NewNotesController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusNoContent)
	})

	t.Run("Fail to delete note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		note := domain.Note{ID: id, UserID: user.ID}

		repoMock := storage.NewMockNotesRepo(ctrl)
		repoMock.EXPECT().Delete(note).Return(fmt.Errorf("error"))

		c := NewNotesController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})
}
