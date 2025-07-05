package api

import (
	"mywall-api/internal/models"
	"net/http"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mywall-api/internal/helpers"
)

type GalleryRequest struct {
	Title    	string `json:"title" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=500"`
	CategoryID  uint   `json:"category_id" binding:"required"`
}

func (s *Server) getGalleries(c *gin.Context) {
	userID := c.GetUint("user_id")
	var galleries []models.Gallery
	s.db.Where("user_id = ?", userID).Find(&galleries)
	helpers.Success(c, "Gallies retrieved successfully", galleries)	
}

func (s *Server) getGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var gallery models.Gallery
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}
	helpers.Success(c, "Gallery retrieved successfully", gallery)	
}

func (s *Server) createGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Parse form data for file upload
	var req GalleryRequest
	var finalImageURL string
	
	// Get form values manually for better error handling
	categoryIDStr := c.PostForm("category_id")
	if categoryIDStr == "" {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"category_id": "Category ID is required",
		})
		return
	}
	
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"category_id": "Category ID must be a valid number",
		})
		return
	}
	
	req.CategoryID = uint(categoryID)
	req.Title = strings.TrimSpace(c.PostForm("title"))
	req.Description = strings.TrimSpace(c.PostForm("description"))
	
	// Validate required fields
	if req.Title == "" {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"title": "Title is required",
		})
		return
	}
	
	// Validate title length
	if len(req.Title) > 100 {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"title": "Title must not exceed 100 characters",
		})
		return
	}
	
	// Validate description length
	if len(req.Description) > 500 {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"description": "Description must not exceed 500 characters",
		})
		return
	}
	
	// Handle file upload - required
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		helpers.BadRequest(c, "Image file is required")
		return
	}
	defer file.Close()
	
	// Validate file type
	if !isValidImageFile(header.Filename) {
		helpers.BadRequest(c, "Invalid file type. Only JPG and PNG files are allowed")
		return
	}
	
	// Validate file size (5MB limit)
	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if header.Size > maxFileSize {
		helpers.BadRequest(c, "File size too large. Maximum allowed size is 5MB")
		return
	}
	
	// Check if user exists
	var user models.User
	if result := s.db.First(&user, userID); result.Error != nil {
		helpers.NotFound(c, "Invalid user")
		return
	}
	
	// Check if category exists
	var category models.Category
	if result := s.db.First(&category, req.CategoryID); result.Error != nil {
		helpers.BadRequest(c, "Invalid category")
		return
	}
	
	// Create directory structure and save file
	filePath, err := saveUploadedFile(file, header)
	if err != nil {
		helpers.InternalServerError(c, "Failed to save image file")
		return
	}
	
	finalImageURL = filePath
	
	// Create gallery record
	gallery := models.Gallery{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    finalImageURL,
		CategoryID:  req.CategoryID,
		UserID:      userID,
	}
	
	if result := s.db.Create(&gallery); result.Error != nil {
		// If database creation fails, clean up the uploaded file
		os.Remove(finalImageURL)
		helpers.InternalServerError(c, "Failed to create gallery")
		return
	}
	
	helpers.Created(c, "Gallery created successfully", gallery)
}

func (s *Server) updateGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var gallery models.Gallery
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&gallery).Error; err != nil {
		helpers.NotFound(c,"Gallery not found")
		return
	}

	var input models.Gallery
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c,err.Error())
		return
	}

	s.db.Model(&gallery).Updates(models.Gallery{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		CategoryID:  input.CategoryID,
	})
	helpers.Success(c, "Gallery updated successfully", gallery)
}

func (s *Server) deleteGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var gallery models.Gallery
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&gallery).Error; err != nil {
		helpers.NotFound(c, "Gallery not found")
		return
	}
	s.db.Delete(&gallery)
	helpers.Success(c, "Gallery deleted", gallery)	
}

// Helper function to validate image file extensions
func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

// Helper function to save uploaded file with structured path /2025/06/05/filename.jpg
func saveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
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

// Helper function to delete file when gallery is deleted
func deleteGalleryFile(imageURL string) error {
	// Only delete if it's a local file (not an external URL)
	if !strings.HasPrefix(imageURL, "http") && imageURL != "" {
		return os.Remove(imageURL)
	}
	return nil
}
