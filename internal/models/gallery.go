package models

import "gorm.io/gorm"

// Gallery represents the gallery item
type Gallery struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	UserID      uint   `json:"user_id"`
}
