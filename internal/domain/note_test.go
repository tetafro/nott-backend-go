package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteValidation(t *testing.T) {
	cases := []struct {
		title string
		note  Note
		err   bool
	}{
		{
			title: "correct note",
			note: Note{
				UserID:    10,
				NotepadID: 20,
				Title:     "x-note",
				Text:      "hello, world",
			},
			err: false,
		},
		{
			title: "note without user",
			note: Note{
				NotepadID: 20,
				Title:     "x-note",
				Text:      "hello, world",
			},
			err: true,
		},
		{
			title: "note without notepad",
			note: Note{
				UserID: 10,
				Title:  "x-note",
				Text:   "hello, world",
			},
			err: true,
		},
		{
			title: "note without title",
			note: Note{
				UserID:    10,
				NotepadID: 20,
				Text:      "hello, world",
			},
			err: true,
		},
		{
			title: "note without text",
			note: Note{
				UserID:    10,
				NotepadID: 20,
				Title:     "x-note",
			},
			err: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			err := tt.note.Validate()
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
