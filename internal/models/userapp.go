package models

import "gorm.io/gorm"

// UserApp represents the user application item
type UserApp struct {
	gorm.Model
	ID          string `json:"id" gorm:"primaryKey"`
	AppID       string `json:"app_id" gorm:"not null;index"`
	UserID      uint    `json:"user_id" gorm:"not null;index"`
	RoleID      string  `json:"role_id" gorm:"not null;index"`
	Status      string  `json:"status" gorm:"not null;index"`
	LastLogin   string  `json:"last_login"`
}