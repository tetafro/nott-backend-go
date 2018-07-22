package folders

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/tetafro/nott-backend-go/internal/database"
	"github.com/tetafro/nott-backend-go/internal/errors"
)

// Repo deals with folders repository.
type Repo interface {
	Get(filter) ([]Folder, error)
	Create(Folder) (Folder, error)
	Update(Folder) (Folder, error)
	Delete(Folder) error
}

// PostgresRepo is a folders repository that uses PostgreSQL as a backend.
type PostgresRepo struct {
	db *gorm.DB
}

// NewPostgresRepo creates new PostgreSQL repository.
func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Get gets folders from repository.
func (r *PostgresRepo) Get(f filter) ([]Folder, error) {
	ff := []Folder{}

	q := r.db
	if f.id != nil {
		q = q.Where("id = ?", *f.id)
	}
	if f.userID != nil {
		q = q.Where("user_id = ?", *f.userID)
	}

	if err := q.Find(&ff).Error; err != nil {
		return nil, fmt.Errorf("query failed with error: %v", err)
	}

	return ff, nil
}

// Create creates folder in repository.
func (r *PostgresRepo) Create(f Folder) (Folder, error) {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		if err = tx.Create(&f).Error; err != nil {
			return fmt.Errorf("query failed with error: %v", err)
		}
		return nil
	})
	if err != nil {
		return Folder{}, err
	}
	return f, nil
}

// Update updates folder in repository.
func (r *PostgresRepo) Update(f Folder) (Folder, error) {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		// Check if folder exists
		err = tx.Select("id").
			Where("id = ? AND user_id = ?", f.ID, f.UserID).
			Find(&Folder{}).
			Error
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("failed to check folder in database: %v", err)
		}

		// NOTE: Save() method doesn't return ErrRecordNotFound, but
		// instead makes INSERT. But this is the only method that updates
		// all fields of the structure (even if they are empty).
		if err = tx.Save(&f).Error; err != nil {
			return fmt.Errorf("query failed with error: %v", err)
		}

		return nil
	})
	if err != nil {
		return Folder{}, err
	}
	return f, nil
}

// Delete deletes folder in repository.
func (r *PostgresRepo) Delete(f Folder) error {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		err = tx.Where("id = ? AND user_id = ?", f.ID, f.UserID).Delete(&Folder{}).Error
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
	id     *uint
	userID *uint
}
