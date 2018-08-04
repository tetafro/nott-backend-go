package auth

import (
	"testing"
)

func TestHash(t *testing.T) {
	t.Run("Hash password and check result", func(t *testing.T) {
		password := "qwerty"
		hash, err := hashPassword(password)
		if err != nil {
			t.Fatalf("Failed to get hash: %v", err)
		}
		if !checkPassword(password, hash) {
			t.Fatalf("Password and hash don't match (but should)")
		}
	})
}
