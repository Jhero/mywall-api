package api

import (
	"mywall-api/internal/models"
	"net/http"
	"fmt"
	"os"
	"strconv"
	"strings"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"errors"
)

type ImageViewRequest struct {
	GalleryID string `json:"gallery_id"`
	Count     int    `json:"count"`
}

func (s *Server) createImageView(ctx *gin.Context) {
	userID := c.GetUint("user_id")
	var req ImageViewRequest
	if req.GalleryID == "" {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"gallery_id": "GalleryID is required",
		})
		return
	}
	if req.Count <= 0 {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"count": "Count must be greater than 0",
		})
		return
	}
	imageView := models.ImageView{
		GalleryID: req.GalleryID,
		Count:     req.Count,
		UserID:    userID,
	}
	if result := s.db.Create(&imageView); result.Error != nil {
		// If database creation fails and we uploaded a file, clean it up
		helpers.InternalServerError(c, "Failed to create image view")
		return
	}	
	helpers.Created(c, "Image view created successfully", imageView)
}

func (s *Server) updateImageView(ctx *gin.Context) {
	userID := c.GetUint("user_id")
	var req ImageViewRequest
	if req.GalleryID == "" {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"gallery_id": "GalleryID is required",
		})
		return
	}
	if req.Count <= 0 {
		helpers.ValidationError(c, "Validation failed", map[string]string{
			"count": "Count must be greater than 0",
		})
		return
	}
	var imageView models.ImageView
	if result := s.db.Where("user_id = ? AND gallery_id = ?", userID, req.GalleryID).First(&imageView); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			helpers.NotFound(c, "Image view not found")
		} else {
			helpers.InternalServerError(c, "Failed to get image view")
		}
		return
	}
	imageView.Count = imageView.Count + req.Count
	if result := s.db.Save(&imageView); result.Error != nil {
		helpers.InternalServerError(c, "Failed to update image view")
		return
	}
	helpers.OK(c, "Image view updated successfully", imageView)
}
