package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// NotesRepo is a notes repository that uses PostgreSQL as a backend.
type NotesRepo struct {
	db *gorm.DB
}

// NewNotesRepo creates new PostgreSQL repository for notes.
func NewNotesRepo(db *gorm.DB) *NotesRepo {
	return &NotesRepo{db: db}
}

// Get gets notes from repository.
func (r *NotesRepo) Get(f storage.NotesFilter) ([]domain.Note, error) {
	n := []domain.Note{}

	q := r.db
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.NotepadID != nil {
		q = q.Where("notepad_id = ?", *f.NotepadID)
	}

	if err := q.Find(&n).Error; err != nil {
		return nil, fmt.Errorf("query failed with error: %v", err)
	}

	return n, nil
}

// Create creates note in repository.
func (r *NotesRepo) Create(n domain.Note) (domain.Note, error) {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		if err = tx.Create(&n).Error; err != nil {
			return fmt.Errorf("query failed with error: %v", err)
		}
		return nil
	})
	if err != nil {
		return domain.Note{}, err
	}
	return n, nil
}

// Update updates note in repository.
func (r *NotesRepo) Update(n domain.Note) (domain.Note, error) {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		// Check if note exists
		err = tx.Select("id").
			Where("id = ? AND user_id = ?", n.ID, n.UserID).
			Find(&domain.Note{}).
			Error
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("failed to check note in database: %v", err)
		}

		// NOTE: Save() method doesn't return ErrRecordNotFound, but
		// instead makes INSERT. But this is the only method that updates
		// all fields of the structure (even if they are empty).
		if err = tx.Save(&n).Error; err != nil {
			return fmt.Errorf("query failed with error: %v", err)
		}

		return nil
	})
	if err != nil {
		return domain.Note{}, err
	}
	return n, nil
}

// Delete deletes note in repository.
func (r *NotesRepo) Delete(n domain.Note) error {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		err = tx.Where("id = ? AND user_id = ?", n.ID, n.UserID).Delete(&domain.Note{}).Error
		if err != nil {
			return fmt.Errorf("query failed with error: %v", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
