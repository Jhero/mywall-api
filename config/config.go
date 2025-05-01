package config

import (
	"os"
	"strconv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	OAuthConfig 		*oauth2.Config
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
		OAuthConfig: 		&oauth2.Config{
			ClientID:     	os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: 	os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  	os.Getenv("OAUTH_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: 		google.Endpoint,
		},
	}
}
