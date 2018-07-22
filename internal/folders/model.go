package folders

// Folder represents a folder that contains notepads.
type Folder struct {
	ID       uint   `json:"id" gorm:"column:id"`
	UserID   uint   `json:"user_id" gorm:"column:user_id"`
	ParentID *uint  `json:"parent_id" gorm:"column:parent_id"`
	Title    string `json:"title" gorm:"column:title"`
}
