package auth

import (
	"regexp"
	"time"

	"github.com/pkg/errors"
)

var regexpEmail = regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]+`)

// User represents a user that used for authenticating.
type User struct {
	ID       int    `json:"id" gorm:"column:id"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"-" gorm:"column:password"`

	// Managed by gorm callbacks
	CreatedAt time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"-" gorm:"column:updated_at"`
}

// Validate validates user.
func (u User) Validate() error {
	if u.Email == "" {
		return errors.New("email cannot be empty")
	}
	if !regexpEmail.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
