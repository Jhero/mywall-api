package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"mywall-api/config"
	"mywall-api/internal/api"
	"mywall-api/internal/database"
	"mywall-api/internal/auth"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables based on environment
	loadEnvFiles()

	// Parse command line flags
	useSqlMigrations := flag.Bool("sql-migrations", false, "Use SQL migrations instead of GORM AutoMigrate")
	migrationsDir := flag.String("migrations-dir", "migrations", "Directory for SQL migrations")
	flag.Parse()

	// Initialize config
	cfg := config.New()

	// Log environment info
	log.Printf("🚀 Starting application in %s mode", cfg.Environment)
	log.Printf("📍 Port: %s, Domain: %s, HTTPS: %t", cfg.Port, cfg.Domain, cfg.UseHTTPS)
	log.Printf("🔧 Debug mode: %t", cfg.Debug)

	// Setup database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected successfully")

	// Run migrations
	if *useSqlMigrations {
		absPath, err := filepath.Abs(*migrationsDir)
		if err != nil {
			log.Fatalf("❌ Failed to get absolute path for migrations: %v", err)
		}
		
		log.Printf("📦 Running SQL migrations from %s", absPath)
		if err := database.MigrateWithSQL(db, absPath); err != nil {
			log.Fatalf("❌ Failed to run SQL migrations: %v", err)
		}
		log.Println("✅ SQL migrations completed successfully")
	} else {
		log.Println("📦 Running GORM AutoMigrate")
		if err := database.Migrate(db); err != nil {
			log.Fatalf("❌ Failed to run migrations: %v", err)
		}
		log.Println("✅ GORM AutoMigrate completed successfully")
	}

	// Initialize auth service with JWT secret
	authService := auth.NewService(db, cfg.JWTSecret)
	log.Println("✅ Auth service initialized")

	// Initialize and start the server
	server := api.NewServer(db, authService)
	
	// Add environment-specific middleware if needed
	if cfg.IsProduction() {
		log.Println("🛡️  Production mode: Security headers enabled")
	} else {
		log.Println("🔍 Development mode: Debug features enabled")
	}

	log.Printf("🌐 Starting server on %s...", cfg.GetServerURL())
	log.Fatal(server.Start(cfg.Port))
}

// loadEnvFiles loads environment files based on current environment
func loadEnvFiles() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Try to load specific environment file first
	envFile := ".env." + env
	if err := godotenv.Load(envFile); err == nil {
		log.Printf("✅ Loaded environment from %s", envFile)
	} else {
		// Fallback to default .env file
		if err := godotenv.Load(); err != nil {
			log.Printf("ℹ️  No .env file found, using system environment variables")
		} else {
			log.Println("✅ Loaded environment from .env")
		}
	}
}