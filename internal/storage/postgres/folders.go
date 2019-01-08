package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// FoldersRepo is a folders repository that uses PostgreSQL as a backend.
type FoldersRepo struct {
	db *gorm.DB
}

// NewFoldersRepo creates new PostgreSQL repository for folders.
func NewFoldersRepo(db *gorm.DB) *FoldersRepo {
	return &FoldersRepo{db: db}
}

// Get gets folders from repository.
func (r *FoldersRepo) Get(f storage.FoldersFilter) ([]domain.Folder, error) {
	ff := []domain.Folder{}

	q := r.db
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.UserID != nil {
		q = q.Where("user_id = ?", *f.UserID)
	}

	if err := q.Find(&ff).Error; err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	return ff, nil
}

// Create creates folder in repository.
func (r *FoldersRepo) Create(f domain.Folder) (domain.Folder, error) {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		if err = tx.Create(&f).Error; err != nil {
			return errors.Wrap(err, "query error")
		}
		return nil
	})
	if err != nil {
		return domain.Folder{}, err
	}
	return f, nil
}

// Update updates folder in repository.
func (r *FoldersRepo) Update(f domain.Folder) (domain.Folder, error) {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		// Check if folder exists
		err = tx.Select("id").
			Where("id = ? AND user_id = ?", f.ID, f.UserID).
			Find(&domain.Folder{}).
			Error
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		if err != nil {
			return errors.Wrap(err, "check folder in database")
		}

		// NOTE: Save() method doesn't return ErrRecordNotFound, but
		// instead makes INSERT. But this is the only method that updates
		// all fields of the structure (even if they are empty).
		if err = tx.Save(&f).Error; err != nil {
			return errors.Wrap(err, "query error")
		}

		return nil
	})
	if err != nil {
		return domain.Folder{}, err
	}
	return f, nil
}

// Delete deletes folder in repository.
func (r *FoldersRepo) Delete(f domain.Folder) error {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		err = tx.Where("id = ? AND user_id = ?", f.ID, f.UserID).Delete(&domain.Folder{}).Error
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
