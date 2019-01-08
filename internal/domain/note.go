package domain

import (
	"fmt"
	"time"
)

// Note represents a note that contains text.
type Note struct {
	ID        int    `json:"id" gorm:"column:id"`
	UserID    int    `json:"user_id" gorm:"column:user_id"`
	NotepadID int    `json:"notepad_id" gorm:"column:notepad_id"`
	Title     string `json:"title" gorm:"column:title"`
	Text      string `json:"text" gorm:"column:text"`
	// HTML field is rendered from Text (which contains markdown) on the fly
	HTML string `json:"html,omitempty" gorm:"-"`

	// Managed by gorm callbacks
	CreatedAt time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"-" gorm:"column:updated_at"`
}

// Validate validates note.
func (n Note) Validate() error {
	if n.UserID == 0 {
		return fmt.Errorf("unknown user")
	}
	if n.NotepadID == 0 {
		return fmt.Errorf("notepad id cannot be empty")
	}
	if n.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	return nil
}
