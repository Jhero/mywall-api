package models

import "gorm.io/gorm"

// User represents the user model
type User struct {
	gorm.Model
	Email        string `json:"email" gorm:"unique"`
	Name         string `json:"name"`
	Password     string `json:"-" gorm:"not null"` // Password hash
	ApiKey       string `json:"api_key" gorm:"unique"`
	Role         string `json:"role" gorm:"default:'user'"`
	IsActive     bool   `json:"is_active" gorm:"default:true"`
}
