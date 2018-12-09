package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// FoldersController handles HTTP API requests.
type FoldersController struct {
	repo storage.FoldersRepo
	log  logrus.FieldLogger
}

// NewFoldersController creates new controller.
func NewFoldersController(repo storage.FoldersRepo, log logrus.FieldLogger) *FoldersController {
	return &FoldersController{repo: repo, log: log}
}

// GetList handles request for getting folders.
func (c *FoldersController) GetList(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)

	folders, err := c.repo.Get(storage.FoldersFilter{UserID: &user.ID})
	if err != nil {
		c.log.Errorf("Failed to get folders: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, folders)
}

// Create handles request for creating folder.
func (c *FoldersController) Create(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)

	var err error

	f := domain.Folder{}
	if err = json.NewDecoder(req.Body).Decode(&f); err != nil {
		badRequest(w, "invalid json")
		return
	}
	f.UserID = user.ID

	if err = f.Validate(); err != nil {
		badRequest(w, "invalid folder: "+err.Error())
		return
	}

	f, err = c.repo.Create(f)
	if err != nil {
		c.log.Errorf("Failed to create folder: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusCreated, f)
}

// GetOne handles request for getting folder by id.
func (c *FoldersController) GetOne(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	folders, err := c.repo.Get(storage.FoldersFilter{ID: &id, UserID: &user.ID})

	if err != nil {
		c.log.Errorf("Failed to get folder: %v", err)
		internalServerError(w)
		return
	}
	if len(folders) == 0 {
		notFound(w)
		return
	}

	respond(w, http.StatusOK, folders[0])
}

// Update handles request for updating folder.
func (c *FoldersController) Update(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	f := domain.Folder{}
	if err = json.NewDecoder(req.Body).Decode(&f); err != nil {
		badRequest(w, "invalid json")
		return
	}
	f.ID = id
	f.UserID = user.ID

	if err = f.Validate(); err != nil {
		badRequest(w, "invalid folder: "+err.Error())
		return
	}

	f, err = c.repo.Update(f)
	if err == domain.ErrNotFound {
		notFound(w)
		return
	}
	if err != nil {
		c.log.Errorf("Failed to update folder: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, f)
}

// Delete handles request for deleting folder.
func (c *FoldersController) Delete(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	id, err := getID(req)
	if err != nil {
		notFound(w)
		return
	}

	f := domain.Folder{ID: id, UserID: user.ID}
	if err = c.repo.Delete(f); err != nil {
		c.log.Errorf("Failed to update folder: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusNoContent, nil)
}
