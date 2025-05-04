package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Port        		string
	JWTSecret        	string
	JWTExpiryHours   	int
	APIKeyHeader     	string
	DBHost           	string
	DBPort           	string
	DBUser           	string
	DBPassword       	string
	DBName           	string	
	DatabaseURL 		string
}

// New creates a new Config with values from environment variables
func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	jwtExpiryHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))

	return &Config{
		Port:        		port,
		JWTSecret:        	os.Getenv("JWT_SECRET"),
		JWTExpiryHours:   	jwtExpiryHours,
		APIKeyHeader:     	os.Getenv("API_KEY_HEADER"),
		DatabaseURL: 		os.Getenv("DATABASE_URL"),
	}
}
