package domain

import (
	"fmt"
	"time"
)

// Folder represents a folder that contains notepads.
type Folder struct {
	ID       int    `json:"id" gorm:"column:id"`
	UserID   int    `json:"user_id" gorm:"column:user_id"`
	ParentID *int   `json:"parent_id" gorm:"column:parent_id"`
	Title    string `json:"title" gorm:"column:title"`

	// Managed by gorm callbacks
	CreatedAt time.Time  `json:"-" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"-" gorm:"column:updated_at"`
}

// Validate validates folder.
func (f Folder) Validate() error {
	if f.UserID == 0 {
		return fmt.Errorf("unknown user")
	}
	if f.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	return nil
}
