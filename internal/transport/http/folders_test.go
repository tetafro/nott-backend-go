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

func TestFoldersController(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	Int := func(n int) *int {
		return &n
	}
	user := auth.User{ID: 1}

	t.Run("Get folders", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		folders := []domain.Folder{
			{ID: 10, UserID: user.ID, ParentID: Int(30), Title: "Folder 10"},
			{ID: 15, UserID: user.ID, ParentID: Int(35), Title: "Folder 15"},
		}

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.FoldersFilter{UserID: &user.ID},
		).Return(folders, nil)

		c := NewFoldersController(repoMock, log)

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
					"user_id": 1,
					"parent_id": 30,
					"title": "Folder 10"
				},
				{
					"id": 15,
					"user_id": 1,
					"parent_id": 35,
					"title": "Folder 15"
				}
			]
		}`)
	})

	t.Run("Fail to get folders", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.FoldersFilter{UserID: &user.ID},
		).Return(nil, fmt.Errorf("error"))

		c := NewFoldersController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)

		c.GetList(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Create folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		folder := domain.Folder{ID: id, UserID: user.ID, ParentID: Int(30), Title: "Folder 10"}

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Create(folder).Return(folder, nil)

		c := NewFoldersController(repoMock, log)

		payload, err := json.Marshal(folder)
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
				"parent_id": 30,
				"title": "Folder 10"
			}
		}`)
	})

	t.Run("Fail to create folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		folder := domain.Folder{ID: id, UserID: user.ID, ParentID: Int(30), Title: "Folder 10"}

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Update(folder).Return(domain.Folder{}, fmt.Errorf("error"))

		c := NewFoldersController(repoMock, log)

		payload, err := json.Marshal(folder)
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

	t.Run("Get folder by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		folders := []domain.Folder{
			{ID: id, UserID: user.ID, ParentID: Int(30), Title: "Folder 10"},
		}

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.FoldersFilter{ID: &id, UserID: &user.ID},
		).Return(folders, nil)

		c := NewFoldersController(repoMock, log)

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
				"parent_id": 30,
				"title": "Folder 10"
			}
		}`)
	})

	t.Run("Fail to get folder by non-existing id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.FoldersFilter{ID: &id, UserID: &user.ID},
		).Return(nil, nil)

		c := NewFoldersController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.GetOne(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	})

	t.Run("Fail to get folder by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Get(
			storage.FoldersFilter{ID: &id, UserID: &user.ID},
		).Return(nil, fmt.Errorf("error"))

		c := NewFoldersController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.GetOne(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	})

	t.Run("Update folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		folder := domain.Folder{ID: id, UserID: user.ID, ParentID: Int(30), Title: "Folder 10"}

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Update(folder).Return(folder, nil)

		c := NewFoldersController(repoMock, log)

		payload, err := json.Marshal(folder)
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
				"parent_id": 30,
				"title": "Folder 10"
			}
		}`)
	})

	t.Run("Fail to update folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		folder := domain.Folder{ID: id, UserID: user.ID, ParentID: Int(30), Title: "Folder 10"}

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Update(folder).Return(domain.Folder{}, fmt.Errorf("error"))

		c := NewFoldersController(repoMock, log)

		payload, err := json.Marshal(folder)
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

	t.Run("Delete folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		folder := domain.Folder{ID: id, UserID: user.ID}

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Delete(folder).Return(nil)

		c := NewFoldersController(repoMock, log)

		url := "/"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req = addUser(req, user)
		req = addID(req, id)

		c.Delete(w, req)

		resp := w.Result()
		assert.Equal(t, resp.StatusCode, http.StatusNoContent)
	})

	t.Run("Fail to delete folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := 10
		folder := domain.Folder{ID: id, UserID: user.ID}

		repoMock := storage.NewMockFoldersRepo(ctrl)
		repoMock.EXPECT().Delete(folder).Return(fmt.Errorf("error"))

		c := NewFoldersController(repoMock, log)

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
