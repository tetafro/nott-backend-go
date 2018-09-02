package notepads

import "testing"

func TestValidation(t *testing.T) {
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
			if tt.err && err == nil {
				t.Fatal("Wanted error, got nil")
			}
			if !tt.err && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}
