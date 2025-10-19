// internal/api/utils.go
package api

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	_ "image/gif"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
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

// SaveAndOptimizeImage saves and optimizes the uploaded image with compression and resize
func SaveAndOptimizeImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image to maximum dimensions (800x800) while maintaining aspect ratio
	maxWidth := 50
	maxHeight := 50
	img = resizeImage(img, maxWidth, maxHeight)

	// Get current date for directory structure
	now := time.Now()
	dateDir := fmt.Sprintf("/%d/%02d/%02d", now.Year(), now.Month(), now.Day())
	
	// Create base upload directory
	baseDir := "uploads"
	fullDir := filepath.Join(baseDir, dateDir)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate unique filename
	ext := strings.ToLower(filepath.Ext(header.Filename))
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(fullDir, uniqueFilename)
	
	// Create output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// Encode with compression based on format
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(outFile, img)
	default:
		// Default to JPEG for other formats
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 80})
	}

	if err != nil {
		os.Remove(filePath) // Clean up on error
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return filePath, nil
}

// resizeImage resizes an image maintaining aspect ratio
func resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate new dimensions maintaining aspect ratio
	var newWidth, newHeight int
	if width > maxWidth || height > maxHeight {
		ratio := float64(width) / float64(height)
		
		if width > height {
			newWidth = maxWidth
			newHeight = int(float64(maxWidth) / ratio)
		} else {
			newHeight = maxHeight
			newWidth = int(float64(maxHeight) * ratio)
		}
		
		// Use Lanczos resampling for better quality
		return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
	}
	
	return img
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