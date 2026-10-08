package models

import (
	"time"
)

type PostInteraction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_post_user" json:"post_id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_post_user;comment:'UserID integration'" json:"user_id"`
	Type      string    `gorm:"not null" json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}