package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"mywall-api/config"
	"mywall-api/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Parse command-line flags
	create := flag.String("create", "", "Create a new migration")
	migrationsDir := flag.String("dir", "migrations", "Directory for migrations")
	flag.Parse()

	// Ensure migrations directory exists
	if err := os.MkdirAll(*migrationsDir, 0755); err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	// Create migration if requested
	if *create != "" {
		filePath, err := database.CreateMigration(*migrationsDir, *create)
		if err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		fmt.Printf("Created migration file: %s\n", filePath)
		return
	}

	// Otherwise, run migrations
	cfg := config.New()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	absPath, err := filepath.Abs(*migrationsDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	if err := database.MigrateWithSQL(db, absPath); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	fmt.Println("All migrations applied successfully")
}
