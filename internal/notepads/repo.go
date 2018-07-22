package notepads

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/tetafro/nott-backend-go/internal/database"
	"github.com/tetafro/nott-backend-go/internal/errors"
)

// Repo deals with notepads repository.
type Repo interface {
	Get(filter) ([]Notepad, error)
	Create(Notepad) (Notepad, error)
	Update(Notepad) (Notepad, error)
	Delete(Notepad) error
}

// PostgresRepo is a notepads repository that uses PostgreSQL as a backend.
type PostgresRepo struct {
	db *gorm.DB
}

// NewPostgresRepo creates new PostgreSQL repository.
func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Get gets notepads from repository.
func (r *PostgresRepo) Get(f filter) ([]Notepad, error) {
	n := []Notepad{}

	q := r.db
	if f.id != nil {
		q = q.Where("id = ?", *f.id)
	}
	if f.userID != nil {
		q = q.Where("user_id = ?", *f.userID)
	}
	if f.folderID != nil {
		q = q.Where("folder_id = ?", *f.folderID)
	}

	if err := q.Find(&n).Error; err != nil {
		return nil, fmt.Errorf("query failed with error: %v", err)
	}

	return n, nil
}

// Create creates notepad in repository.
func (r *PostgresRepo) Create(n Notepad) (Notepad, error) {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		if err = tx.Create(&n).Error; err != nil {
			return fmt.Errorf("query failed with error: %v", err)
		}
		return nil
	})
	if err != nil {
		return Notepad{}, err
	}
	return n, nil
}

// Update updates notepad in repository.
func (r *PostgresRepo) Update(n Notepad) (Notepad, error) {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		// Check if notepad exists
		err = tx.Select("id").
			Where("id = ? AND user_id = ?", n.ID, n.UserID).
			Find(&Notepad{}).
			Error
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("failed to check notepad in database: %v", err)
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
		return Notepad{}, err
	}
	return n, nil
}

// Delete deletes notepad in repository.
func (r *PostgresRepo) Delete(n Notepad) error {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		err = tx.Where("id = ? AND user_id = ?", n.ID, n.UserID).Delete(&Notepad{}).Error
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
	id       *uint
	userID   *uint
	folderID *uint
}
