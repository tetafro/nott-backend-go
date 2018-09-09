package auth

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/tetafro/nott-backend-go/internal/database"
	"github.com/tetafro/nott-backend-go/internal/errors"
)

// UsersRepo deals with users repository.
type UsersRepo interface {
	GetByEmail(email string) (User, error)
	GetByToken(token string) (User, error)
	Update(User) (User, error)
}

// UsersPostgresRepo is a users repository that uses PostgreSQL as a backend.
type UsersPostgresRepo struct {
	db *gorm.DB
}

// NewUsersPostgresRepo creates new PostgreSQL repository.
func NewUsersPostgresRepo(db *gorm.DB) *UsersPostgresRepo {
	return &UsersPostgresRepo{db: db}
}

// GetByEmail gets user by his email from repository.
func (r *UsersPostgresRepo) GetByEmail(email string) (User, error) {
	var u User

	q := r.db.Where("email = ?", email)
	err := q.Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("query failed with error: %v", err)
	}

	return u, nil
}

// GetByToken gets user by his token from repository.
func (r *UsersPostgresRepo) GetByToken(token string) (User, error) {
	var u User

	q := r.db.Joins(`JOIN token ON token.user_id = "user".id`).
		Where("token.string = ? AND token.created_at + token.ttl * INTERVAL '1 second' > NOW()", token)
	err := q.Find(&u).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("query failed with error: %v", err)
	}

	return u, nil
}

// Update updates user in repository.
func (r *UsersPostgresRepo) Update(u User) (User, error) {
	err := database.Transact(r.db, func(tx *gorm.DB) (err error) {
		// Check if user exists
		err = tx.Select("id").
			Where("id = ?", u.ID).
			Find(&User{}).
			Error
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
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
		return User{}, err
	}
	return u, nil
}
