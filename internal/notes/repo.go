package notes

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/tetafro/nott-backend-go/internal/database"
	"github.com/tetafro/nott-backend-go/internal/errors"
)

// Repo deals with notes repository.
type Repo interface {
	Get(filter) ([]Note, error)
	Create(Note) (Note, error)
	Update(Note) (Note, error)
	Delete(Note) error
}

// PostgresRepo is a notes repository that uses PostgreSQL as a backend.
type PostgresRepo struct {
	db *gorm.DB
}

// NewPostgresRepo creates new PostgreSQL repository.
func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Get gets notes from repository.
func (r *PostgresRepo) Get(f filter) ([]Note, error) {
	n := []Note{}

	q := r.db
	if f.id != nil {
		q = q.Where("id = ?", *f.id)
	}
	if f.userID != nil {
		q = q.Where("user_id = ?", *f.userID)
	}
	if f.notepadID != nil {
		q = q.Where("notepad_id = ?", *f.notepadID)
	}

	if err := q.Find(&n).Error; err != nil {
		return nil, fmt.Errorf("query failed with error: %v", err)
	}

	return n, nil
}

// Create creates note in repository.
func (r *PostgresRepo) Create(n Note) (Note, error) {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		if err = tx.Create(&n).Error; err != nil {
			return fmt.Errorf("query failed with error: %v", err)
		}
		return nil
	})
	if err != nil {
		return Note{}, err
	}
	return n, nil
}

// Update updates note in repository.
func (r *PostgresRepo) Update(n Note) (Note, error) {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		// Check if note exists
		err = tx.Select("id").
			Where("id = ? AND user_id = ?", n.ID, n.UserID).
			Find(&Note{}).
			Error
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
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
		return Note{}, err
	}
	return n, nil
}

// Delete deletes note in repository.
func (r *PostgresRepo) Delete(n Note) error {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		err = tx.Where("id = ? AND user_id = ?", n.ID, n.UserID).Delete(&Note{}).Error
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

type filter struct {
	id        *uint
	userID    *uint
	notepadID *uint
}
