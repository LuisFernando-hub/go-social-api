package models

import "time"

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`
	UserName string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}