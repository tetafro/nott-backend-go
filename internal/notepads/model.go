package notepads

// Notepad represents a notepad that contains notes.
type Notepad struct {
	ID       uint   `json:"id" gorm:"column:id"`
	UserID   uint   `json:"user_id" gorm:"column:user_id"`
	FolderID *uint  `json:"folder_id" gorm:"column:folder_id"`
	Title    string `json:"title" gorm:"column:title"`
}
