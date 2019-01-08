package postgres

import (
	"math/rand"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/tetafro/nott-backend-go/internal/auth"
)

const (
	// tokenLen is a length of a token string.
	tokenLen = 8

	// defaultTokenTTL is a duration of time in seconds
	// before token will become expired (1 year).
	defaultTokenTTL = 60 * 60 * 24 * 365
)

// TokensRepo is a tokens repository that uses PostgreSQL as a backend.
type TokensRepo struct {
	db *gorm.DB
}

// NewTokensRepo creates new PostgreSQL repository.
func NewTokensRepo(db *gorm.DB) *TokensRepo {
	return &TokensRepo{db: db}
}

// Create creates a new token in repository.
func (r *TokensRepo) Create(t auth.Token) (auth.Token, error) {
	if t.UserID == 0 {
		return auth.Token{}, errors.New("user id is empty")
	}
	if t.String == "" {
		t.String = randString(tokenLen)
	}
	if t.TTL == 0 {
		t.TTL = defaultTokenTTL
	}

	if err := r.db.Create(&t).Error; err != nil {
		return auth.Token{}, errors.Wrap(err, "query error")
	}

	return t, nil
}

// Delete deletes token from repository.
func (r *TokensRepo) Delete(t auth.Token) error {
	if err := r.db.Delete(&t).Error; err != nil {
		return errors.Wrap(err, "query error")
	}
	return nil
}

func randString(n int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
