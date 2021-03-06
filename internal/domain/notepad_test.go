package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotepadValidation(t *testing.T) {
	cases := []struct {
		title   string
		notepad Notepad
		err     bool
	}{
		{
			title: "correct notepad",
			notepad: Notepad{
				UserID:   10,
				FolderID: 20,
				Title:    "x-notepad",
			},
			err: false,
		},
		{
			title: "notepad without user",
			notepad: Notepad{
				FolderID: 20,
				Title:    "x-notepad",
			},
			err: true,
		},
		{
			title: "notepad without folder",
			notepad: Notepad{
				UserID: 10,
				Title:  "x-notepad",
			},
			err: true,
		},
		{
			title: "notepad without title",
			notepad: Notepad{
				UserID:   10,
				FolderID: 20,
			},
			err: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			err := tt.notepad.Validate()
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
