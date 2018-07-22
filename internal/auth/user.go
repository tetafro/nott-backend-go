package auth

// User represents a user that used for authenticating.
type User struct {
	ID       uint   `json:"id" gorm:"column:id"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"-" gorm:"column:password"`
}
