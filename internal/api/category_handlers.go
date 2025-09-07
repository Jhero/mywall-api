package api

import (
	"mywall-api/internal/models"
	"net/http"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mywall-api/internal/helpers"
)

type CategoryRequest struct {
	Name    	string `json:"name" binding:"required,max=50"`
}

/*
func (s *Server) getCategories(c *gin.Context) {
	userID := c.GetUint("user_id")
	var categories []models.Category
	s.db.Where("user_id = ?", userID).Find(&categories)
	c.JSON(http.StatusOK, categories)
}
*/

func (s *Server) getCategories(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Get query parameters for filtering
	name := c.Query("name")
	
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
	
	// Get total count for pagination
	var total int64
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		helpers.NotFound(c, "Failed to count categories")
		return
	}
	
	// Get categories with pagination
	var categories []models.Category
	if err := query.
		Order("created_at DESC").
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
	}
	
	helpers.Success(c, "Categories retrieved successfully", response)
}

func (s *Server) getCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var category models.Category
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (s *Server) createCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				switch e.Field() {
				case "Name":
					if e.Tag() == "required" {
						errorMessages["name"] = "Name is required" 
					}
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Check if user exists
	var user models.User
	if result := s.db.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	category := models.Category{
		Name:       	req.Name,
		UserID:      	userID,
	}
	if result := s.db.Create(&category); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (s *Server) updateCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var category models.Category
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.db.Model(&category).Updates(models.Category{
		Name:       input.Name,
	})

	c.JSON(http.StatusOK, category)
}

func (s *Server) deleteCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var category models.Category
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	s.db.Delete(&category)
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
