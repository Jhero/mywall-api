package models

import "gorm.io/gorm"

// Gallery represents the gallery item
type Role struct {
	gorm.Model
	ID 			string `json:"ID" gorm:"unique"`
	Name       	string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	UserID     	uint   `json:"user_id"`
}
