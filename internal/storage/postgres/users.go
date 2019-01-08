package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/domain"
)

// UsersRepo is a users repository that uses PostgreSQL as a backend.
type UsersRepo struct {
	db *gorm.DB
}

// NewUsersRepo creates new PostgreSQL repository.
func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

// GetByID gets user by his email from repository.
func (r *UsersRepo) GetByID(id int) (auth.User, error) {
	var u auth.User

	err := r.db.Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		return auth.User{}, domain.ErrNotFound
	}
	if err != nil {
		return auth.User{}, errors.Wrap(err, "query error")
	}

	return u, nil
}

// GetByEmail gets user by his email from repository.
func (r *UsersRepo) GetByEmail(email string) (auth.User, error) {
	var u auth.User

	q := r.db.Where("email = ?", email)
	err := q.Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		return auth.User{}, domain.ErrNotFound
	}
	if err != nil {
		return auth.User{}, errors.Wrap(err, "query error")
	}

	return u, nil
}

// Create creates user in repository.
func (r *UsersRepo) Create(u auth.User) (auth.User, error) {
	q := r.db.Create(&u)
	if err := q.Error; err != nil {
		return auth.User{}, errors.Wrap(err, "query error")
	}
	q.Scan(&u)
	return u, nil
}

// Update updates user in repository.
func (r *UsersRepo) Update(u auth.User) (auth.User, error) {
	err := transact(r.db, func(tx *gorm.DB) (err error) {
		// Check if user exists
		err = tx.Select("id").
			Where("id = ?", u.ID).
			Find(&auth.User{}).
			Error
		if err == gorm.ErrRecordNotFound {
			return domain.ErrNotFound
		}
		if err != nil {
			return errors.Wrap(err, "check user in database")
		}

		// NOTE: Save() method doesn't return ErrRecordNotFound, but
		// instead makes INSERT. But this is the only method that updates
		// all fields of the structure (even if they are empty).
		if err = tx.Save(&u).Error; err != nil {
			return errors.Wrap(err, "query error")
		}

		return nil
	})
	if err != nil {
		return auth.User{}, err
	}
	return u, nil
}
