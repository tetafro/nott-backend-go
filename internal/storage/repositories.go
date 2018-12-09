package storage

import (
	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/domain"
)

// UsersRepo deals with users repository.
type UsersRepo interface {
	GetByEmail(email string) (auth.User, error)
	GetByToken(token string) (auth.User, error)
	Update(auth.User) (auth.User, error)
}

// TokensRepo deals with tokens repository.
type TokensRepo interface {
	Create(auth.Token) (auth.Token, error)
	Delete(auth.Token) error
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
	ID     *uint
	UserID *uint
}

// NotepadsFilter is a filter for searching notepads in repository.
type NotepadsFilter struct {
	ID       *uint
	UserID   *uint
	FolderID *uint
}

// NotesFilter is a filter for searching notes in repository.
type NotesFilter struct {
	ID        *uint
	UserID    *uint
	NotepadID *uint
}
