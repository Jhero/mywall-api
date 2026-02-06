package models

import (
    "time"
    "gorm.io/gorm"
)

type Notification struct {
    ID        string         `json:"id" gorm:"type:uuid;primaryKey"`
    Title     string         `json:"title"`
    Body      string         `json:"body"`
    Metadata  string         `json:"metadata"`
    Type      string         `json:"type" gorm:"not null"`
    IsRead    int            `json:"is_read" gorm:"default:0"`
    UserID    uint           `json:"user_id"`

    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
