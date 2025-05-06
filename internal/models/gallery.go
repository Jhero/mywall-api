package models

import "gorm.io/gorm"

// Gallery represents the gallery item
type Gallery struct {
	gorm.Model
	Title       string `json:"title" gorm:"unique"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url" gorm:"not null"`
	CategoryID  uint   `json:"category_id" gorm:"not null"`
	UserID      uint   `json:"user_id"`
}
