package models

import (
	"time"
)

type Post struct {
	ID uint `gorm:"primaryKey" json:"id"`
	UserID uint `json:"user_id"`

	Content string `json:"content"`
	LikeCount uint `json:"like_count"`
	DeslikeCount uint `json:"deslike_count"`

	Interactions []PostInteraction `gorm:"foreignKey:PostID;references:ID" json:"interactions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}