package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserValidation(t *testing.T) {
	cases := []struct {
		title string
		user  User
		err   bool
	}{
		{
			title: "correct user",
			user: User{
				ID:    10,
				Email: "bob@example.com",
			},
			err: false,
		},
		{
			title: "user without email",
			user: User{
				ID: 10,
			},
			err: true,
		},
		{
			title: "user with invalid email #1",
			user: User{
				ID:    10,
				Email: "bob",
			},
			err: true,
		},
		{
			title: "user with invalid email #2",
			user: User{
				ID:    10,
				Email: "bob@example",
			},
			err: true,
		},
		{
			title: "user with invalid email #2",
			user: User{
				ID:    10,
				Email: "bob$@example.com",
			},
			err: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
