package notepads

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

// GetList handles request for getting notepads.
func (c *Controller) GetList(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)

	notepads, err := c.repo.Get(filter{userID: &user.ID})
	if err != nil {
		c.log.Errorf("Failed to get notepads: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(notepads).
		Write(w)
}

// GetOne handles request for getting notepad by id.
func (c *Controller) GetOne(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	notepads, err := c.repo.Get(filter{id: &id, userID: &user.ID})
	if err != nil {
		c.log.Errorf("Failed to get notepad: %v", err)
		response.InternalServerError().Write(w)
		return
	}
	if len(notepads) == 0 {
		response.NotFound().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(notepads[0]).
		Write(w)
}

// Create handles request for creating notepad.
func (c *Controller) Create(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)

	var err error

	n := Notepad{}
	if err = json.NewDecoder(req.Body).Decode(&n); err != nil {
		response.New().
			WithStatus(http.StatusBadRequest).
			WithError(response.Error("Invalid JSON")).
			Write(w)
		return
	}
	n.UserID = user.ID

	n, err = c.repo.Create(n)
	if err != nil {
		c.log.Errorf("Failed to create notepad: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusCreated).
		WithData(n).
		Write(w)
}

// Update handles request for updating notepad.
func (c *Controller) Update(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	n := Notepad{}
	if err = json.NewDecoder(req.Body).Decode(&n); err != nil {
		response.New().
			WithStatus(http.StatusBadRequest).
			WithError(response.Error("Invalid JSON")).
			Write(w)
		return
	}
	n.ID = id
	n.UserID = user.ID

	n, err = c.repo.Update(n)
	if err == errors.ErrNotFound {
		response.NotFound().Write(w)
		return
	}
	if err != nil {
		c.log.Errorf("Failed to update notepad: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(n).
		Write(w)
}

// Delete handles request for deleting notepad.
func (c *Controller) Delete(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	n := Notepad{ID: id, UserID: user.ID}

	if err = c.repo.Delete(n); err != nil {
		c.log.Errorf("Failed to update notepad: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusNoContent).
		WithData(n).
		Write(w)
}
