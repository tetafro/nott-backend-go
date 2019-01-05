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

func TestNotepadsController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	user := auth.User{ID: 1}

	t.Run("Get notepads", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notepads := []domain.Notepad{
			{ID: 10, UserID: 20, FolderID: 30, Title: "Notepad 10"},
			{ID: 15, UserID: 25, FolderID: 35, Title: "Notepad 15"},
		}

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotepadsFilter{UserID: &user.ID},
		).Return(notepads, nil)

		c := NewNotepadsController(repoMock, log)

		url := "/"
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

	t.Run("Fail to get notepads", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotepadsFilter{UserID: &user.ID},
		).Return(nil, fmt.Errorf("error"))

		c := NewNotepadsController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)

		c.GetList(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Create notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := domain.Notepad{ID: id, UserID: user.ID, FolderID: 30, Title: "Notepad 10"}

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Create(notepad).Return(notepad, nil)

		c := NewNotepadsController(repoMock, log)

		payload, err := json.Marshal(notepad)
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
				"folder_id": 30,
				"title": "Notepad 10"
			}
		}`)
	})

	t.Run("Fail to create notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := domain.Notepad{ID: id, UserID: user.ID, FolderID: 30, Title: "Notepad 10"}

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Update(notepad).Return(domain.Notepad{}, fmt.Errorf("error"))

		c := NewNotepadsController(repoMock, log)

		payload, err := json.Marshal(notepad)
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

	t.Run("Get notepad by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepads := []domain.Notepad{
			{ID: id, UserID: 20, FolderID: 30, Title: "Notepad 10"},
		}

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotepadsFilter{ID: &id, UserID: &user.ID},
		).Return(notepads, nil)

		c := NewNotepadsController(repoMock, log)

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
				"user_id": 20,
				"folder_id": 30,
				"title": "Notepad 10"
			}
		}`)
	})

	t.Run("Fail to get notepad by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotepadsFilter{ID: &id, UserID: &user.ID},
		).Return(nil, fmt.Errorf("error"))

		c := NewNotepadsController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.GetOne(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Fail to get notepad by non-existing id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.NotepadsFilter{ID: &id, UserID: &user.ID},
		).Return(nil, nil)

		c := NewNotepadsController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.GetOne(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	})

	t.Run("Update notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := domain.Notepad{ID: id, UserID: user.ID, FolderID: 30, Title: "Notepad 10"}

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Update(notepad).Return(notepad, nil)

		c := NewNotepadsController(repoMock, log)

		payload, err := json.Marshal(notepad)
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
				"folder_id": 30,
				"title": "Notepad 10"
			}
		}`)
	})

	t.Run("Fail to update notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := domain.Notepad{ID: id, UserID: user.ID, FolderID: 30, Title: "Notepad 10"}

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Update(notepad).Return(domain.Notepad{}, fmt.Errorf("error"))

		c := NewNotepadsController(repoMock, log)

		payload, err := json.Marshal(notepad)
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

	t.Run("Delete notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := domain.Notepad{ID: id, UserID: user.ID}

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Delete(notepad).Return(nil)

		c := NewNotepadsController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusNoContent)
	})

	t.Run("Fail to delete notepad", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		notepad := domain.Notepad{ID: id, UserID: user.ID}

		repoMock := storage.NewMockNotepadsRepo(ctrl)
		repoMock.EXPECT().Delete(notepad).Return(fmt.Errorf("error"))

		c := NewNotepadsController(repoMock, log)

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
