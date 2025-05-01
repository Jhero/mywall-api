package main

import (
	"flag"
	"log"
	// "os"
	"path/filepath"

	"mywall-api/config"
	"mywall-api/internal/api"
	"mywall-api/internal/database"
	"mywall-api/internal/auth"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Parse command line flags
	useSqlMigrations := flag.Bool("sql-migrations", false, "Use SQL migrations instead of GORM AutoMigrate")
	migrationsDir := flag.String("migrations-dir", "migrations", "Directory for SQL migrations")
	flag.Parse()

	// Initialize config
	cfg := config.New()

	// Setup database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if *useSqlMigrations {
		absPath, err := filepath.Abs(*migrationsDir)
		if err != nil {
			log.Fatalf("Failed to get absolute path for migrations: %v", err)
		}
		
		log.Printf("Running SQL migrations from %s", absPath)
		if err := database.MigrateWithSQL(db, absPath); err != nil {
			log.Fatalf("Failed to run SQL migrations: %v", err)
		}
		log.Println("SQL migrations completed successfully")
	} else {
		log.Println("Running GORM AutoMigrate")
		if err := database.Migrate(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("GORM AutoMigrate completed successfully")
	}

	// Initialize auth service with JWT secret
	authService := auth.NewService(db, cfg.JWTSecret)

	// Initialize and start the server
	server := api.NewServer(db, authService)
	log.Printf("Starting server on port %s...", cfg.Port)
	log.Fatal(server.Start(cfg.Port))
}