package models

import "gorm.io/gorm"

// Gallery represents the gallery item
type Category struct {
	gorm.Model
	Name       string `json:"name" gorm:"unique"`
	UserID     uint    `json:"user_id"`
	ImageURL   string `json:"image_url" gorm:"not null"`
}
