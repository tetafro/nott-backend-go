package notes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/russross/blackfriday"
	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/errors"
	"github.com/tetafro/nott-backend-go/internal/http/request"
	"github.com/tetafro/nott-backend-go/internal/http/response"
)

// Controller handles HTTP API requests.
type Controller struct {
	repo Repo
	log  logrus.FieldLogger
}

// NewController creates new controller.
func NewController(repo Repo, log logrus.FieldLogger) *Controller {
	return &Controller{repo: repo, log: log}
}

// GetList handles request for getting notes.
func (c *Controller) GetList(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	notepadID := req.URL.Query().Get("notepad_id")

	f := filter{userID: &user.ID}

	if notepadID != "" {
		nid, err := strconv.Atoi(notepadID)
		if err != nil {
			response.New().
				WithStatus(http.StatusBadRequest).
				WithError(response.Error("Notepad ID must be an integer number")).
				Write(w)
			return
		}
		unid := uint(nid)
		f.notepadID = &unid
	}

	notes, err := c.repo.Get(f)
	if err != nil {
		c.log.Errorf("Failed to get notes: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(notes).
		Write(w)
}

// GetOne handles request for getting note by id.
func (c *Controller) GetOne(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	notes, err := c.repo.Get(filter{id: &id, userID: &user.ID})
	if err != nil {
		c.log.Errorf("Failed to get note: %v", err)
		response.InternalServerError().Write(w)
		return
	}
	if len(notes) == 0 {
		response.NotFound().Write(w)
		return
	}

	// Render markdown to HTML
	n := notes[0]
	n.HTML = string(blackfriday.Run([]byte(n.Text)))

	response.New().
		WithStatus(http.StatusOK).
		WithData(n).
		Write(w)
}

// Create handles request for creating note.
func (c *Controller) Create(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)

	var err error

	n := Note{}
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
		c.log.Errorf("Failed to create note: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusCreated).
		WithData(n).
		Write(w)
}

// Update handles request for updating note.
func (c *Controller) Update(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	n := Note{}
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
		c.log.Errorf("Failed to update note: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(n).
		Write(w)
}

// Delete handles request for deleting note.
func (c *Controller) Delete(w http.ResponseWriter, req *http.Request) {
	user := auth.GetUser(req)
	id, err := request.GetUintPathParam(req, "id")
	if err != nil {
		response.NotFound().Write(w)
		return
	}

	n := Note{ID: id, UserID: user.ID}

	if err = c.repo.Delete(n); err != nil {
		c.log.Errorf("Failed to update note: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusNoContent).
		WithData(n).
		Write(w)
}
