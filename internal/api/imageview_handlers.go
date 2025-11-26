package api

import (
	"mywall-api/internal/models"
	"mywall-api/internal/helpers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"errors"
)

type ImageViewRequest struct {
	GalleryID string `json:"gallery_id"`
	Count     int    `json:"count"`
}

func (s *Server) createImageView(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req ImageViewRequest
	if req.GalleryID == "" {
		helpers.ValidationError(ctx, "Validation failed", map[string]string{
			"gallery_id": "GalleryID is required",
		})
		return
	}
	if req.Count <= 0 {
		helpers.ValidationError(ctx, "Validation failed", map[string]string{
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
		helpers.InternalServerError(ctx, "Failed to create image view")
		return
	}	
	helpers.Created(ctx, "Image view created successfully", imageView)
}

func (s *Server) updateImageView(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req ImageViewRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(ctx, "Validation failed", map[string]string{
			"request": "Invalid request body",
		})
		return
	}
	
	if req.GalleryID == "" {
		helpers.ValidationError(ctx, "Validation failed", map[string]string{
			"gallery_id": "GalleryID is required",
		})
		return
	}
	if req.Count <= 0 {
		helpers.ValidationError(ctx, "Validation failed", map[string]string{
			"count": "Count must be greater than 0",
		})
		return
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		helpers.InternalServerError(ctx, "Failed to start transaction")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			helpers.InternalServerError(ctx, "Transaction failed")
		}
	}()

	var imageView models.ImageView
	if result := s.db.Where("user_id = ? AND gallery_id = ?", userID, req.GalleryID).First(&imageView); result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			helpers.NotFound(ctx, "Image view not found")
		} else {
			helpers.InternalServerError(ctx, "Failed to get image view")
		}
		return
	}
	imageView.Count = imageView.Count + req.Count
	if result := s.db.Save(&imageView); result.Error != nil {
		tx.Rollback()
		helpers.InternalServerError(ctx, "Failed to update image view")
		return
	}

	if result := tx.Commit(); result.Error != nil {
		tx.Rollback()
		helpers.InternalServerError(ctx, "Failed to commit transaction")
		return
	}
	helpers.Success(ctx, "Image view updated successfully", imageView)
}
