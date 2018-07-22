package folders

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/errors"
	"github.com/tetafro/nott-backend-go/internal/httpx/request"
	"github.com/tetafro/nott-backend-go/internal/httpx/response"
)

// Controller handles HTTP API requests.
type Controller struct {
	repo Repo
	log  *logrus.Logger
}

// NewController creates new controller.
func NewController(repo Repo, log *logrus.Logger) *Controller {
	return &Controller{repo: repo, log: log}
}

// GetList handles request for getting folders.
func (c *Controller) GetList(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)

	folders, err := c.repo.Get(filter{userID: &user.ID})
	if err != nil {
		c.log.Errorf("Failed to get folders: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(folders).
		Write(w)
}

// Create handles request for creating folder.
func (c *Controller) Create(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)

	var err error

	f := Folder{}
	if err = json.NewDecoder(req.Body).Decode(&f); err != nil {
		response.New().
			WithStatus(http.StatusBadRequest).
			WithError(response.Error("Invalid JSON")).
			Write(w)
		return
	}
	f.UserID = user.ID

	f, err = c.repo.Create(f)
	if err != nil {
		c.log.Errorf("Failed to create folder: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusCreated).
		WithData(f).
		Write(w)
}

// GetOne handles request for getting folder by id.
func (c *Controller) GetOne(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	folders, err := c.repo.Get(filter{id: &id, userID: &user.ID})

	if err != nil {
		c.log.Errorf("Failed to get folder: %v", err)
		response.InternalServerError().Write(w)
		return
	}
	if len(folders) == 0 {
		response.NotFound().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(folders[0]).
		Write(w)
}

// Update handles request for updating folder.
func (c *Controller) Update(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	f := Folder{}
	if err = json.NewDecoder(req.Body).Decode(&f); err != nil {
		response.New().
			WithStatus(http.StatusBadRequest).
			WithError(response.Error("Invalid JSON")).
			Write(w)
		return
	}
	f.ID = id
	f.UserID = user.ID

	f, err = c.repo.Update(f)
	if err == errors.ErrNotFound {
		response.NotFound().Write(w)
		return
	}
	if err != nil {
		c.log.Errorf("Failed to update folder: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(f).
		Write(w)
}

// Delete handles request for deleting folder.
func (c *Controller) Delete(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	f := Folder{ID: id, UserID: user.ID}

	if err = c.repo.Delete(f); err != nil {
		c.log.Errorf("Failed to update folder: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusNoContent).
		WithData(f).
		Write(w)
}
