package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// NotepadsRepo is a notepads repository that uses PostgreSQL as a backend.
type NotepadsRepo struct {
	db *gorm.DB
}

// NewNotepadsRepo creates new PostgreSQL repository for notepads.
func NewNotepadsRepo(db *gorm.DB) *NotepadsRepo {
	return &NotepadsRepo{db: db}
}

// Get gets notepads from repository.
func (r *NotepadsRepo) Get(f storage.NotepadsFilter) ([]domain.Notepad, error) {
	n := []domain.Notepad{}

	q := r.db
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}
	if f.FolderID != nil {
		q = q.Where("folder_id = ?", *f.FolderID)
	}

	if err := q.Find(&n).Error; err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	return n, nil
}

// Create creates notepad in repository.
func (r *NotepadsRepo) Create(n domain.Notepad) (domain.Notepad, error) {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		if err = tx.Create(&n).Error; err != nil {
			return errors.Wrap(err, "query error")
		}
		return nil
	})
	if err != nil {
		return domain.Notepad{}, err
	}
	return n, nil
}

// Update updates notepad in repository.
func (r *NotepadsRepo) Update(n domain.Notepad) (domain.Notepad, error) {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		// Check if notepad exists
		err = tx.Select("id").
			Where("id = ? AND user_id = ?", n.ID, n.UserID).
			Find(&domain.Notepad{}).
			Error
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		if err != nil {
			return errors.Wrap(err, "check notepad in database")
		}

		// NOTE: Save() method doesn't return ErrRecordNotFound, but
		// instead makes INSERT. But this is the only method that updates
		// all fields of the structure (even if they are empty).
		if err = tx.Save(&n).Error; err != nil {
			return errors.Wrap(err, "query error")
		}

		return nil
	})
	if err != nil {
		return domain.Notepad{}, err
	}
	return n, nil
}

// Delete deletes notepad in repository.
func (r *NotepadsRepo) Delete(n domain.Notepad) error {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		err = tx.Where("id = ? AND user_id = ?", n.ID, n.UserID).Delete(&domain.Notepad{}).Error
		if err != nil {
			return errors.Wrap(err, "query error")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
