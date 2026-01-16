package models

import "gorm.io/gorm"

// Gallery represents the gallery item
type Rbac struct {
	gorm.Model
	Permission  string `json:"permission"`
	MenuID      string `json:"menu_id"`
	UserID      uint    `json:"user_id"`
	RoleID 		string `json:"role_id" gorm:"not null;index"`

}

// TableName specifies the table name for Rbac
func (Rbac) TableName() string {
	return "rbacs"
}