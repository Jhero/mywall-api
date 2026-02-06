package models

import "gorm.io/gorm"

// UserApp represents the user application item
type Setting struct {
	gorm.Model
	ID          string `json:"id" gorm:"primaryKey"`
	AppID       string `json:"app_id" gorm:"not null;index"`
	UserID      uint    `json:"user_id" gorm:"not null;index"`
	Key         string  `json:"key" gorm:"not null;index"`
	Value       string  `json:"value" gorm:"not null;index"`
}