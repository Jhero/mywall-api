package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Environment      string
	Port             string
	JWTSecret        string
	JWTExpiryHours   int
	APIKeyHeader     string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DatabaseURL      string
	Domain           string
	UseHTTPS         bool
	Debug            bool
}

// New creates a new Config with values from environment variables
func New() *Config {
	// Get environment (default to development)
	env := getEnv("APP_ENV", "development")
	
	// Get port with default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get JWT expiry hours
	jwtExpiryHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if jwtExpiryHours == 0 {
		jwtExpiryHours = 24 // default 24 hours
	}

	return &Config{
		Environment:    env,
		Port:           port,
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiryHours: jwtExpiryHours,
		APIKeyHeader:   os.Getenv("API_KEY_HEADER"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		Domain:         getDomain(env),
		UseHTTPS:       getUseHTTPS(env),
		Debug:          getDebug(env),
	}
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Get domain based on environment
func getDomain(env string) string {
	if env == "production" {
		return getEnv("DOMAIN", "myjovan.site")
	}
	return "localhost"
}

// Determine if HTTPS should be used
func getUseHTTPS(env string) bool {
	if env == "production" {
		useHTTPS, _ := strconv.ParseBool(getEnv("USE_HTTPS", "true"))
		return useHTTPS
	}
	return false
}

// Determine if debug mode should be enabled
func getDebug(env string) bool {
	if env == "development" {
		debug, _ := strconv.ParseBool(getEnv("DEBUG", "true"))
		return debug
	}
	return false
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetServerURL returns the appropriate server URL based on environment
func (c *Config) GetServerURL() string {
	if c.IsProduction() && c.UseHTTPS {
		return "https://" + c.Domain
	}
	return "http://" + c.Domain + ":" + c.Port
}