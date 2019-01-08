package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJub3R0Iiwic3ViIjoiMTAifQ.q4tyUIUMs7H3ikFuxhmlNBENL9Hqh8XVt0DQNmyTvHU" // nolint
	secret := "qwerty"
	userID := 10

	t.Run("Issue token", func(t *testing.T) {
		tokener := NewJWTokener(secret)
		tokener.ttl = 0

		token, err := tokener.Issue(User{ID: userID})
		assert.NoError(t, err)
		assert.Equal(t, accessToken, token.AccessToken)
	})

	t.Run("Parse token", func(t *testing.T) {
		tokener := NewJWTokener(secret)

		id, err := tokener.Parse(accessToken)
		assert.NoError(t, err)
		assert.Equal(t, userID, id)
	})

	t.Run("Fail to parse token", func(t *testing.T) {
		tokener := NewJWTokener(secret)

		_, err := tokener.Parse("malformed.token")
		assert.Error(t, err)
	})
}
