package models

import "gorm.io/gorm"

// App represents the application item
type App struct {
	gorm.Model
	Permission  string `json:"permission"`
	MenuID      string `json:"menu_id"`
	UserID      uint    `json:"user_id"`
	RoleID 		string `json:"role_id" gorm:"not null;index"`

}