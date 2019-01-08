package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// NotepadsController handles HTTP API requests.
type NotepadsController struct {
	repo storage.NotepadsRepo
	log  logrus.FieldLogger
}

// NewNotepadsController creates new controller.
func NewNotepadsController(repo storage.NotepadsRepo, log logrus.FieldLogger) *NotepadsController {
	return &NotepadsController{repo: repo, log: log}
}

// GetList handles request for getting notepads.
func (c *NotepadsController) GetList(w http.ResponseWriter, req *http.Request) {
	userID := getUserID(req)

	notepads, err := c.repo.Get(storage.NotepadsFilter{UserID: &userID})
	if err != nil {
		c.log.Errorf("Failed to get notepads: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, notepads)
}

// GetOne handles request for getting notepad by id.
func (c *NotepadsController) GetOne(w http.ResponseWriter, req *http.Request) {
	userID := getUserID(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	notepads, err := c.repo.Get(storage.NotepadsFilter{ID: &id, UserID: &userID})
	if err != nil {
		c.log.Errorf("Failed to get notepad: %v", err)
		internalServerError(w)
		return
	}
	if len(notepads) == 0 {
		notFound(w)
		return
	}

	respond(w, http.StatusOK, notepads[0])
}

// Create handles request for creating notepad.
func (c *NotepadsController) Create(w http.ResponseWriter, req *http.Request) {
	userID := getUserID(req)

	var err error

	n := domain.Notepad{}
	if err = json.NewDecoder(req.Body).Decode(&n); err != nil {
		badRequest(w, "invalid json")
		return
	}
	n.UserID = userID

	if err = n.Validate(); err != nil {
		badRequest(w, "invalid notepad"+err.Error())
		return
	}

	n, err = c.repo.Create(n)
	if err != nil {
		c.log.Errorf("Failed to create notepad: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusCreated, n)
}

// Update handles request for updating notepad.
func (c *NotepadsController) Update(w http.ResponseWriter, req *http.Request) {
	userID := getUserID(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	n := domain.Notepad{}
	if err = json.NewDecoder(req.Body).Decode(&n); err != nil {
		badRequest(w, "invalid json")
		return
	}
	n.ID = id
	n.UserID = userID

	if err = n.Validate(); err != nil {
		badRequest(w, "invalid notepad"+err.Error())
		return
	}

	n, err = c.repo.Update(n)
	if err == domain.ErrNotFound {
		notFound(w)
		return
	}
	if err != nil {
		c.log.Errorf("Failed to update notepad: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, n)
}

// Delete handles request for deleting notepad.
func (c *NotepadsController) Delete(w http.ResponseWriter, req *http.Request) {
	userID := getUserID(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	n := domain.Notepad{ID: id, UserID: userID}
	if err = c.repo.Delete(n); err != nil {
		c.log.Errorf("Failed to update notepad: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusNoContent, nil)
}
