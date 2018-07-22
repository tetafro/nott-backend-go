package auth

import "time"

// Token is used for user authentication.
type Token struct {
	ID      uint      `json:"-"  gorm:"column:id"`
	UserID  uint      `json:"-" gorm:"column:user_id"`
	String  string    `json:"string" gorm:"column:string"`
	TTL     int       `json:"ttl" gorm:"column:ttl"` // seconds
	Created time.Time `json:"created" gorm:"column:created"`
}
