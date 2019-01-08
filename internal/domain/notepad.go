package domain

import (
	"fmt"
	"time"
)

// Notepad represents a notepad that contains notes.
type Notepad struct {
	ID       int    `json:"id" gorm:"column:id"`
	UserID   int    `json:"user_id" gorm:"column:user_id"`
	FolderID int    `json:"folder_id" gorm:"column:folder_id"`
	Title    string `json:"title" gorm:"column:title"`

	// Managed by gorm callbacks
	CreatedAt time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"-" gorm:"column:updated_at"`
}

// Validate validates notepad.
func (n Notepad) Validate() error {
	if n.UserID == 0 {
		return fmt.Errorf("unknown user")
	}
	if n.FolderID == 0 {
		return fmt.Errorf("folder id cannot be empty")
	}
	if n.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	return nil
}
