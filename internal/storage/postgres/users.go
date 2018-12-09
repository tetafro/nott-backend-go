package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"

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

// GetByEmail gets user by his email from repository.
func (r *UsersRepo) GetByEmail(email string) (auth.User, error) {
	var u auth.User

	q := r.db.Where("email = ?", email)
	err := q.Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		return auth.User{}, domain.ErrNotFound
	}
	if err != nil {
		return auth.User{}, fmt.Errorf("query failed with error: %v", err)
	}

	return u, nil
}

// GetByToken gets auth.user by his token from repository.
func (r *UsersRepo) GetByToken(token string) (auth.User, error) {
	var u auth.User

	q := r.db.Joins(`JOIN token ON token.user_id = "user".id`).
		Where("token.string = ? AND token.created_at + token.ttl * INTERVAL '1 second' > NOW()", token)
	err := q.Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		return auth.User{}, domain.ErrNotFound
	}
	if err != nil {
		return auth.User{}, fmt.Errorf("query failed with error: %v", err)
	}

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
			return fmt.Errorf("failed to check user in database: %v", err)
		}

		// NOTE: Save() method doesn't return ErrRecordNotFound, but
		// instead makes INSERT. But this is the only method that updates
		// all fields of the structure (even if they are empty).
		if err = tx.Save(&u).Error; err != nil {
			return fmt.Errorf("query failed with error: %v", err)
		}

		return nil
	})
	if err != nil {
		return auth.User{}, err
	}
	return u, nil
}