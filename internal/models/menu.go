package models

import "gorm.io/gorm"

// Gallery represents the gallery item
type Menu struct {
	gorm.Model
	ID 			string `json:"ID" gorm:"unique"`
	Path 		string `json:"path" gorm:"unique"`
	UserID      uint    `json:"user_id"`
}
