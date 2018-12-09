package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderValidation(t *testing.T) {
	parentID := uint(10)

	cases := []struct {
		title  string
		folder Folder
		err    bool
	}{
		{
			title: "correct folder",
			folder: Folder{
				UserID:   10,
				ParentID: &parentID,
				Title:    "x-folder",
			},
			err: false,
		},
		{
			title: "folder without user",
			folder: Folder{
				ParentID: &parentID,
				Title:    "x-folder",
			},
			err: true,
		},
		{
			title: "folder without parent",
			folder: Folder{
				UserID: 10,
				Title:  "x-folder",
			},
			err: false,
		},
		{
			title: "folder without title",
			folder: Folder{
				UserID:   10,
				ParentID: &parentID,
			},
			err: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.title, func(t *testing.T) {
			err := tt.folder.Validate()
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
