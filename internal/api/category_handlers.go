package api

import (
	"mywall-api/internal/models"
	"math"
	"strconv"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"mywall-api/internal/helpers"
)

type CategoryRequest struct {
	Name    	string `json:"name" binding:"required,max=50"`
}

func (s *Server) getCategories(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Get query parameters for filtering
	name := c.Query("name")
	
	// Get sorting parameters
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	
	// Get pagination parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	
	// Convert pagination parameters
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		limitInt = 10
	}
	
	offset := (pageInt - 1) * limitInt
	
	// Build query with base condition
	query := s.db.Model(&models.Category{}).Where("user_id = ?", userID)
	
	if name != "" {
		// Case-insensitive search using LOWER for better compatibility
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%")
	}
	
	// Validate and build sort order
	var orderBy string
	switch sortBy {
	case "name":
		if sortOrder == "asc" {
			orderBy = "name ASC"
		} else {
			orderBy = "name DESC"
		}
	case "created_at":
		if sortOrder == "asc" {
			orderBy = "created_at ASC"
		} else {
			orderBy = "created_at DESC"
		}
	default:
		orderBy = "created_at DESC" // Default fallback
	}
	
	// Get total count for pagination
	var total int64
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		helpers.NotFound(c, "Failed to count categories")
		return
	}
	
	// Get categories with pagination and sorting
	var categories []models.Category
	if err := query.
		Order(orderBy).
		Limit(limitInt).
		Offset(offset).
		Find(&categories).Error; err != nil {
		helpers.NotFound(c, "Failed to retrieve categories")
		return
	}
	
	// Calculate pagination metadata
	totalPages := int(math.Ceil(float64(total) / float64(limitInt)))
	hasNext := pageInt < totalPages
	hasPrev := pageInt > 1
	
	// Response with metadata
	response := gin.H{
		"data": categories,
		"pagination": gin.H{
			"current_page":   pageInt,
			"total_pages":    totalPages,
			"total_items":    total,
			"items_per_page": limitInt,
			"has_next":       hasNext,
			"has_previous":   hasPrev,
		},
		"filters": gin.H{
			"name":       name,
		},
		"sorting": gin.H{
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	}
	
	helpers.Success(c, "Categories retrieved successfully", response)
}

func (s *Server) getCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var category models.Category
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		helpers.NotFound(c, "Category not found")
		return
	}
	helpers.Success(c, "Category retrieved successfully", category)
}

func (s *Server) createCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CategoryRequest
	var finalImageURL string

	// Get name from form data
	req.Name = strings.TrimSpace(c.PostForm("name"))
	
	// Validate name field
	if req.Name == "" {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"name": "Name is required",
		})
		return
	}
	
	if len(req.Name) > 50 {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"name": "Name must not exceed 50 characters",
		})
		return
	}

	// Check if user exists
	var user models.User
	if result := s.db.First(&user, userID); result.Error != nil {
		helpers.BadRequest(c, "Invalid user")
		return
	}

	// Handle image upload
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		helpers.BadRequest(c, "Image file is required")
		return
	}
	defer file.Close()
	
	// Validate file type using shared utility function
	if !IsValidImageFile(header.Filename) {
		helpers.BadRequest(c, "Invalid file type. Only JPG, JPEG, and PNG files are allowed")
		return
	}
	
	// Validate file size (5MB limit)
	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if header.Size > maxFileSize {
		helpers.BadRequest(c, "File size too large. Maximum allowed size is 5MB")
		return
	}

	// Save the uploaded file using shared utility function
	filePath, err := SaveUploadedFile(file, header)
	if err != nil {
		helpers.InternalServerError(c, "Failed to save image file")
		return
	}
	finalImageURL = filePath

	category := models.Category{
		Name:     req.Name,
		UserID:   userID,
		ImageURL: finalImageURL,
	}
	
	if result := s.db.Create(&category); result.Error != nil {
		// If database creation fails and we uploaded a file, clean it up
		os.Remove(finalImageURL)
		helpers.InternalServerError(c, "Failed to create category")
		return
	}
	
	helpers.Created(c, "Category created successfully", category)
}

func (s *Server) updateCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var category models.Category
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		helpers.NotFound(c, "Category not found")
		return
	}

	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c, "Invalid request data")
		return
	}

	s.db.Model(&category).Updates(models.Category{
		Name: input.Name,
	})

	helpers.Success(c, "Category updated successfully", category)
}

func (s *Server) deleteCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var category models.Category
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		helpers.NotFound(c, "Category not found")
		return
	}

	// Delete associated image file if it exists
	if category.ImageURL != "" {
		if err := os.Remove(category.ImageURL); err != nil {
			// Log the error but don't fail the deletion
			// You might want to use your logging system here
			// fmt.Printf("Warning: Failed to delete image file %s: %v\n", category.ImageURL, err)
		}
	}
	
	s.db.Delete(&category)
	helpers.Success(c, "Category deleted successfully", nil)
}