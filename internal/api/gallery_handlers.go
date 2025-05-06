package api

import (
	"mywall-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mywall-api/internal/helpers"
)

type GalleryRequest struct {
	Title    	string `json:"title" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=500"`
	ImageURL   	string `json:"image_url" binding:"required,url"`
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
	var req GalleryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				switch e.Field() {
					case "CategoryID":
						if e.Tag() == "required" {
							errorMessages["category_id"] = "Category ID is required" 
						}
					case "Title":
						if e.Tag() == "required" {
							errorMessages["title"] = "Title is required" 
						}
					case "ImageURL":
						if e.Tag() == "required" {
							errorMessages["image_url"] = "Image URL is required" 
						}
				}
			}
			helpers.ValidationError(c,"Validation failed", errorMessages)	
			return
		}
		helpers.BadRequest(c,"Invalid request data")	
		return
	}

	// Check if user exists
	var user models.User
	if result := s.db.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	// Check if category exists
	var category models.Category
	if result := s.db.First(&category, req.CategoryID); result.Error != nil {
		// log.Printf("Looking up category with ID: %d", req.CategoryID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
		return
	}
	
	gallery := models.Gallery{
		Title:       	req.Title,
		Description: 	req.Description,
		ImageURL:    	req.ImageURL,
		CategoryID:		req.CategoryID,
		UserID:      	userID,
		// Set other fields as needed
	}
	if result := s.db.Create(&gallery); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gallery"})
		return
	}
	c.JSON(http.StatusCreated, gallery)
}

func (s *Server) updateGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var gallery models.Gallery
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}

	var input models.Gallery
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.db.Model(&gallery).Updates(models.Gallery{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		CategoryID:  input.CategoryID,
	})

	c.JSON(http.StatusOK, gallery)
}

func (s *Server) deleteGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var gallery models.Gallery
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}
	s.db.Delete(&gallery)
	c.JSON(http.StatusOK, gin.H{"message": "Gallery deleted"})
}
