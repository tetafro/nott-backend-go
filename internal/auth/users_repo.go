package auth

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// UsersRepo deals with users repository.
type UsersRepo interface {
	GetByEmail(email string) (User, error)
	GetByToken(token string) (User, error)
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
