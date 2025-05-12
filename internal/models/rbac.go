package models

import "gorm.io/gorm"

// Gallery represents the gallery item
type Rbac struct {
	gorm.Model
	Permission  string `json:"path"`
	MenuID      string `json:"menu_id"`
	UserID      uint   `json:"user_id"`
}
