package notes

// Note represents a note that contains text.
type Note struct {
	ID        uint   `json:"id" gorm:"column:id"`
	UserID    uint   `json:"user_id" gorm:"column:user_id"`
	NotepadID *uint  `json:"notepad_id" gorm:"column:notepad_id"`
	Title     string `json:"title" gorm:"column:title"`
	Text      string `json:"text" gorm:"column:text"`
	HTML      string `json:"html" gorm:"-"`
}
