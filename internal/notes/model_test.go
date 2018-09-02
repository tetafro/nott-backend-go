package notes

import "testing"

func TestValidation(t *testing.T) {
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
			if tt.err && err == nil {
				t.Fatal("Wanted error, got nil")
			}
			if !tt.err && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}
