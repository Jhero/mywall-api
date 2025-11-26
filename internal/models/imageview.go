package models

import "gorm.io/gorm"

// Gallery represents the gallery item
type ImageView struct {
	gorm.Model
	GalleryID  	string `json:"gallery_id" gorm:"not null"`
	Count     	int    `json:"count"`
	UserID     	uint   `json:"user_id"`
}
