package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/markdown"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// NotesController handles HTTP API requests.
type NotesController struct {
	repo storage.NotesRepo
	log  logrus.FieldLogger
}

// NewNotesController creates new controller.
func NewNotesController(repo storage.NotesRepo, log logrus.FieldLogger) *NotesController {
	return &NotesController{repo: repo, log: log}
}

// GetList handles request for getting notes.
func (c *NotesController) GetList(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	notepadID := req.URL.Query().Get("notepad_id")

	f := storage.NotesFilter{UserID: &user.ID}

	if notepadID != "" {
		nid, err := strconv.Atoi(notepadID)
		if err != nil {
			badRequest(w, "Notepad ID must be an integer number")
			return
		}
		f.NotepadID = &nid
	}

	notes, err := c.repo.Get(f)
	if err != nil {
		c.log.Errorf("Failed to get notes: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, notes)
}

// GetOne handles request for getting note by id.
func (c *NotesController) GetOne(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	notes, err := c.repo.Get(storage.NotesFilter{ID: &id, UserID: &user.ID})
	if err != nil {
		c.log.Errorf("Failed to get note: %v", err)
		internalServerError(w)
		return
	}
	if len(notes) == 0 {
		notFound(w)
		return
	}
	n := notes[0]

	// Render markdown to HTML
	n.HTML = markdown.Render(n.Text)

	respond(w, http.StatusOK, n)
}

// Create handles request for creating note.
func (c *NotesController) Create(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)

	var err error

	n := domain.Note{}
	if err = json.NewDecoder(req.Body).Decode(&n); err != nil {
		badRequest(w, "invalid json")
		return
	}
	n.UserID = user.ID

	if err = n.Validate(); err != nil {
		badRequest(w, "invalid note"+err.Error())
		return
	}

	n, err = c.repo.Create(n)
	if err != nil {
		c.log.Errorf("Failed to create note: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusCreated, n)
}

// Update handles request for updating note.
func (c *NotesController) Update(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	n := domain.Note{}
	if err = json.NewDecoder(req.Body).Decode(&n); err != nil {
		badRequest(w, "invalid json")
		return
	}
	n.ID = id
	n.UserID = user.ID

	if err = n.Validate(); err != nil {
		badRequest(w, "invalid note"+err.Error())
		return
	}

	n, err = c.repo.Update(n)
	if err == domain.ErrNotFound {
		notFound(w)
		return
	}
	if err != nil {
		c.log.Errorf("Failed to update note: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, n)
}

// Delete handles request for deleting note.
func (c *NotesController) Delete(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	n := domain.Note{ID: id, UserID: user.ID}
	if err = c.repo.Delete(n); err != nil {
		c.log.Errorf("Failed to update note: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusNoContent, nil)
}
