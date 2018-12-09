package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	t.Run("Hash password and check result", func(t *testing.T) {
		password := "qwerty"
		hash, err := HashPassword(password)
		assert.NoError(t, err)

		match := CheckPassword(password, hash)
		assert.True(t, match)
	})
}
