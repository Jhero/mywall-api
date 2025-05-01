package api

import (
	"mywall-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getGalleries(c *gin.Context) {
	userID := c.GetUint("user_id")
	var galleries []models.Gallery
	s.db.Where("user_id = ?", userID).Find(&galleries)
	c.JSON(http.StatusOK, galleries)
}

func (s *Server) getGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var gallery models.Gallery
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&gallery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gallery not found"})
		return
	}
	c.JSON(http.StatusOK, gallery)
}

func (s *Server) createGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	var gallery models.Gallery
	if err := c.ShouldBindJSON(&gallery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	gallery.UserID = userID
	s.db.Create(&gallery)
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
