package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// tokenLen is a length of a token string.
	tokenLen = 8

	// defaultTokenTTL is a duration of time in seconds
	// before token will become expired (1 year).
	defaultTokenTTL = 60 * 60 * 24 * 365
)

// TokensRepo deals with tokens repository.
type TokensRepo interface {
	Create(Token) (Token, error)
	Delete(Token) error
}

// TokensPostgresRepo is a tokens repository that uses PostgreSQL as a backend.
type TokensPostgresRepo struct {
	db *gorm.DB
}

// NewTokensPostgresRepo creates new PostgreSQL repository.
func NewTokensPostgresRepo(db *gorm.DB) *TokensPostgresRepo {
	return &TokensPostgresRepo{db: db}
}

// Create creates a new token in repository.
func (r *TokensPostgresRepo) Create(t Token) (Token, error) {
	if t.UserID == 0 {
		return Token{}, fmt.Errorf("user id is empty")
	}
	if t.String == "" {
		t.String = randString(tokenLen)
	}
	if t.TTL == 0 {
		t.TTL = defaultTokenTTL
	}
	t.Created = time.Now().UTC()

	if err := r.db.Create(&t).Error; err != nil {
		return Token{}, fmt.Errorf("query failed with error: %v", err)
	}

	return t, nil
}

// Delete deletes token from repository.
func (r *TokensPostgresRepo) Delete(t Token) error {
	if err := r.db.Delete(&t).Error; err != nil {
		return fmt.Errorf("query failed with error: %v", err)
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
