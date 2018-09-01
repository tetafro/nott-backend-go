package folders

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

	t.Run("Get folders", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		folders := []Folder{
			{ID: 10, UserID: user.ID, ParentID: Uint(30), Title: "Folder 10"},
			{ID: 15, UserID: user.ID, ParentID: Uint(35), Title: "Folder 15"},
		}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Get(
			filter{userID: &user.ID},
		).Return(folders, nil)

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

	t.Run("Failed to get folders", func(t *testing.T) {
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
			t.Fatalf("Expected status code %d, but got %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("Create folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		folder := Folder{ID: id, UserID: user.ID, ParentID: Uint(30), Title: "Folder 10"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Create(folder).Return(folder, nil)

		c := NewController(repoMock, log)

		payload, err := json.Marshal(folder)
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
				"parent_id": 30,
				"title": "Folder 10"
			}
		}`)
	})

	t.Run("Failed to create folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		folder := Folder{ID: id, UserID: user.ID, ParentID: Uint(30), Title: "Folder 10"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(folder).Return(Folder{}, fmt.Errorf("error"))

		c := NewController(repoMock, log)

		payload, err := json.Marshal(folder)
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

	t.Run("Get folder by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		folders := []Folder{
			{ID: id, UserID: user.ID, ParentID: Uint(30), Title: "Folder 10"},
		}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Get(
			filter{id: &id, userID: &user.ID},
		).Return(folders, nil)

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
				"parent_id": 30,
				"title": "Folder 10"
			}
		}`)
	})

	t.Run("Failed to get folder by id", func(t *testing.T) {
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
			t.Fatalf("Expected status code %d, but got %d", http.StatusInternalServerError, resp.StatusCode)
		}
	})

	t.Run("Update folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		folder := Folder{ID: id, UserID: user.ID, ParentID: Uint(30), Title: "Folder 10"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(folder).Return(folder, nil)

		c := NewController(repoMock, log)

		payload, err := json.Marshal(folder)
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
				"parent_id": 30,
				"title": "Folder 10"
			}
		}`)
	})

	t.Run("Failed to update folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		folder := Folder{ID: id, UserID: user.ID, ParentID: Uint(30), Title: "Folder 10"}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Update(folder).Return(Folder{}, fmt.Errorf("error"))

		c := NewController(repoMock, log)

		payload, err := json.Marshal(folder)
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

	t.Run("Delete folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		folder := Folder{ID: id, UserID: user.ID}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Delete(folder).Return(nil)

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

	t.Run("Failed to delete folder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := uint(10)
		folder := Folder{ID: id, UserID: user.ID}

		repoMock := NewMockRepo(ctrl)
		repoMock.EXPECT().Delete(folder).Return(fmt.Errorf("error"))

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
