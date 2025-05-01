package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null;unique"`
	AppliedAt time.Time `gorm:"not null"` // Removed the default:CURRENT_TIMESTAMP
}

// MigrateWithSQL runs SQL migrations from the migrations directory
func MigrateWithSQL(db *gorm.DB, migrationsDir string) error {
	// Create migrations table manually instead of using AutoMigrate
	err := createMigrationsTable(db)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedMigrations := make(map[string]bool)
	var migrations []Migration
	if err := db.Find(&migrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	for _, m := range migrations {
		appliedMigrations[m.Name] = true
	}

	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort migration files
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Apply migrations
	for _, fileName := range migrationFiles {
		// Skip if already applied
		if appliedMigrations[fileName] {
			fmt.Printf("Migration %s already applied, skipping\n", fileName)
			continue
		}

		// Read and apply migration
		filePath := filepath.Join(migrationsDir, fileName)
		sqlBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		sql := string(sqlBytes)
		fmt.Printf("Applying migration: %s\n", fileName)

		// Extract the Up section and split into statements
		upSQL := extractUpSection(sql)
		statements := splitSQLStatements(upSQL)

		// Begin transaction
		tx := db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to start transaction for migration %s: %w", fileName, tx.Error)
		}

		// Execute each statement individually
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			
			if err := tx.Exec(stmt).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to apply migration statement in %s: %w\nStatement: %s", fileName, err, stmt)
			}
		}

		// Record migration
		if err := tx.Create(&Migration{Name: fileName, AppliedAt: time.Now()}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", fileName, err)
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", fileName, err)
		}

		fmt.Printf("Successfully applied migration: %s\n", fileName)
	}

	return nil
}

// extractUpSection extracts the SQL statements from the "Up" section of a migration file
func extractUpSection(sql string) string {
	// Find "-- Up" and "-- Down" markers
	upIndex := strings.Index(sql, "-- Up")
	downIndex := strings.Index(sql, "-- Down")
	
	if upIndex == -1 {
		// No explicit sections, use the whole file
		return sql
	}
	
	if downIndex == -1 {
		// No down section, use everything from "-- Up"
		return strings.TrimSpace(sql[upIndex+len("-- Up"):])
	}
	
	// Extract just the Up section
	return strings.TrimSpace(sql[upIndex+len("-- Up"):downIndex])
}

// splitSQLStatements splits a SQL string into individual statements
func splitSQLStatements(sql string) []string {
	// Simple implementation - split by semicolons
	// This will handle most cases but might not handle complex SQL with embedded semicolons in quotes, etc.
	statements := strings.Split(sql, ";")
	
	// Remove any empty statements
	var result []string
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) != "" {
			result = append(result, stmt)
		}
	}
	
	return result
}

// createMigrationsTable creates the migrations table with proper SQL
func createMigrationsTable(db *gorm.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS migrations (
		id bigint unsigned AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		applied_at datetime(3) NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT uni_migrations_name UNIQUE (name)
	)
	`
	return db.Exec(sql).Error
}

// CreateMigration creates a new migration file
func CreateMigration(migrationsDir, name string) (string, error) {
	// Ensure migrations directory exists
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Format filename with timestamp
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filePath := filepath.Join(migrationsDir, fileName)

	// Create file with template
	template := `-- Migration: ` + name + `
-- Created at: ` + time.Now().Format(time.RFC3339) + `
-- Up

-- Write your up migration here

-- Down
-- Uncomment if you want to use down migrations

-- Write your down migration here
`

	if err := ioutil.WriteFile(filePath, []byte(template), 0644); err != nil {
		return "", fmt.Errorf("failed to create migration file: %w", err)
	}

	return filePath, nil
}