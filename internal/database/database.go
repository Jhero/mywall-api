package database

import (
	"mywall-api/internal/models"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect establishes a connection to the database
func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Gallery{}, &models.User{})
}