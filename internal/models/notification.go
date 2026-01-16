package models

import "gorm.io/gorm"

// Notification represents the notification item
type Notification struct {
	gorm.Model
	ID          string `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"unique"`
	Body        string `json:"body"`
	Metadata    string `json:"metadata"`
	Type        string `json:"type" gorm:"not null"`
	IsRead      bool   `json:"is_read" gorm:"default:false"`
	UserID      uint    `json:"user_id"`
}
