package api

import (
	"mywall-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryRequest struct {
	Name    	string `json:"name" binding:"required,max=50"`
}

func (s *Server) getCategories(c *gin.Context) {
	userID := c.GetUint("user_id")
	var categories []models.Category
	s.db.Where("user_id = ?", userID).Find(&categories)
	c.JSON(http.StatusOK, categories)
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
