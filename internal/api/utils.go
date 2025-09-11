// internal/api/utils.go
package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Helper function to save uploaded file with structured path /2025/06/05/filename.jpg
func SaveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Get current date for directory structure
	now := time.Now()
	dateDir := fmt.Sprintf("/%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	
	// Create base upload directory
	baseDir := "uploads" // You can configure this path
	fullDir := filepath.Join(baseDir, dateDir)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}
	
	// Generate unique filename to avoid conflicts
	ext := strings.ToLower(filepath.Ext(header.Filename))
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// Full file path
	filePath := filepath.Join(fullDir, uniqueFilename)
	
	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()
	
	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		// Clean up on error
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	
	return filePath, nil
}

// Helper function to validate image file extensions
func IsValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

// Helper function to validate image file extensions with more formats
func IsValidImageFileExtended(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}