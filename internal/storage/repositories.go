package storage

import (
	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/domain"
)

// UsersRepo deals with users repository.
type UsersRepo interface {
	GetByID(id int) (auth.User, error)
	GetByEmail(email string) (auth.User, error)
	Create(auth.User) (auth.User, error)
	Update(auth.User) (auth.User, error)
}

// FoldersRepo deals with folders repository.
type FoldersRepo interface {
	Get(FoldersFilter) ([]domain.Folder, error)
	Create(domain.Folder) (domain.Folder, error)
	Update(domain.Folder) (domain.Folder, error)
	Delete(domain.Folder) error
}

// NotepadsRepo deals with notepads repository.
type NotepadsRepo interface {
	Get(NotepadsFilter) ([]domain.Notepad, error)
	Create(domain.Notepad) (domain.Notepad, error)
	Update(domain.Notepad) (domain.Notepad, error)
	Delete(domain.Notepad) error
}

// NotesRepo deals with notes repository.
type NotesRepo interface {
	Get(NotesFilter) ([]domain.Note, error)
	Create(domain.Note) (domain.Note, error)
	Update(domain.Note) (domain.Note, error)
	Delete(domain.Note) error
}

// FoldersFilter is a filter for searching foldres in repository.
type FoldersFilter struct {
	ID     *int
	UserID *int
}

// NotepadsFilter is a filter for searching notepads in repository.
type NotepadsFilter struct {
	ID       *int
	UserID   *int
	FolderID *int
}

// NotesFilter is a filter for searching notes in repository.
type NotesFilter struct {
	ID        *int
	UserID    *int
	NotepadID *int
}
