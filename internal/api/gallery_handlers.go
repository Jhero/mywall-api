package api

import (
	"mywall-api/internal/models"
	// "net/http"
	"fmt"
	"os"
	"strconv"
	"strings"
	"math"

	"github.com/gin-gonic/gin"
	"mywall-api/internal/helpers"
	"gorm.io/gorm"
	"math/rand"
	"errors"
	
	// "log"
)
type GalleryRequest struct {
	Title    	string `json:"title" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=500"`
	CategoryID  uint   `json:"category_id" binding:"required"`
}

func (s *Server) getGalleries(c *gin.Context) {
    userID := c.GetUint("user_id")

    categoryID := c.Query("category_id")
	// log.Println("categoryID:", categoryID)
    title := c.Query("title")

    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "10")

    pageInt, err := strconv.Atoi(page)
    if err != nil || pageInt < 1 {
        pageInt = 1
    }

    limitInt, err := strconv.Atoi(limit)
    if err != nil || limitInt < 1 || limitInt > 100 {
        limitInt = 10
    }

    offset := (pageInt - 1) * limitInt

    // Base query
    baseQuery := s.db.Model(&models.Gallery{}).Where("user_id = ?", userID)

    if categoryID != "" {
        baseQuery = baseQuery.Where("category_id = ?", categoryID)
    }

    if title != "" {
        baseQuery = baseQuery.Where("LOWER(title) LIKE LOWER(?)", "%"+title+"%")
    }

    // Count query (fresh session, no limit/offset)
    var total int64
    if err := s.db.Model(&models.Gallery{}).
        Where("user_id = ?", userID).
        Scopes(func(db *gorm.DB) *gorm.DB {
            if categoryID != "" {
                db = db.Where("category_id = ?", categoryID)
            }
            if title != "" {
                db = db.Where("LOWER(title) LIKE LOWER(?)", "%"+title+"%")
            }
            return db
        }).
        Count(&total).Error; err != nil {
        helpers.NotFound(c, "Failed to count galleries")
        return
    }

    // Get paginated data
    var galleries []models.Gallery
    if err := baseQuery.
        Order("created_at DESC").
        Limit(limitInt).
        Offset(offset).
        Find(&galleries).Error; err != nil {
        helpers.NotFound(c, "Failed to retrieve galleries")
        return
    }

	if len(galleries) == 0 {
		helpers.NotContent(c, "No galleries found")
		return
	}

    // Pagination metadata
    totalPages := int(math.Ceil(float64(total) / float64(limitInt)))
    if totalPages == 0 {
        totalPages = 1
    }

    response := gin.H{
        "data": galleries,
        "pagination": gin.H{
            "current_page":   pageInt,
            "total_pages":    totalPages,
            "total_items":    total,
            "items_per_page": limitInt,
            "has_next":       pageInt < totalPages,
            "has_previous":   pageInt > 1,
        },
        "filters": gin.H{
            "category_id": categoryID,
            "title":       title,
        },
    }

    helpers.Success(c, "Galleries retrieved successfully", response)
}


func (s *Server) getGallery(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var gallery models.Gallery
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&gallery).Error; err != nil {
		helpers.NotFound(c, "Gallery not found")
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
	if !IsValidImageFile(header.Filename) {
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
	filePath, err := SaveUploadedFile(file, header)
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

    // Panggil NotificationHandlers
	notifHandler := &NotificationHandlers{db: s.db}
	_ = notifHandler.CreateNotificationDirect(
		userID,
		"Gallery Created",
		"A new gallery has been created",
		"gallery",
		map[string]interface{}{"title": gallery.Title},
	)


	// h := &NotificationHandlers{}
	// h.createNotification(s.db, userID, "Gallery Created", "A new gallery has been created", "gallery", map[string]string{
	// 	"title": gallery.Title,
	// })
	
	if result := s.db.Create(&gallery); result.Error != nil {
		// If database creation fails, clean up the uploaded file
		os.Remove(finalImageURL)
		helpers.InternalServerError(c, "Failed to create gallery")
		return
	}
	
	// Broadcast ke semua connected clients
	BroadcastNewGallery(map[string]interface{}{
		"ID":          gallery.ID,
		"title":       gallery.Title,
		"description": gallery.Description,
		"image_url":   gallery.ImageURL,
		"category_id": gallery.CategoryID,
		"created_at":  gallery.CreatedAt,
		"updated_at":  gallery.UpdatedAt,
	})

	helpers.Created(c, "Gallery created successfully", gallery)
}

func (s *Server) updateGallery(c *gin.Context) {
	// 1. Validasi user_id dari context
	userID := c.GetUint("user_id")
	if userID == 0 {
		helpers.Unauthorized(c, "Invalid user")
		return
	}

	// 2. Validasi dan konversi id parameter
	idStr := c.Param("id")
	if idStr == "" {
		helpers.BadRequest(c, "ID parameter is required")
		return
	}

	// Konversi string ID ke uint dengan validasi
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.BadRequest(c, "Invalid ID format")
		return
	}

	// 3. Cari gallery dengan validasi yang lebih robust
	var gallery models.Gallery
	if err := s.db.Where("id = ? AND user_id = ?", uint(id), userID).First(&gallery).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.NotFound(c, "Gallery not found")
		} else {
			helpers.InternalServerError(c, "Database error")
		}
		return
	}

	// 4. Handle form-data untuk file upload dan data lainnya
	var input struct {
		Title       string `form:"title" binding:"required,min=1,max=100"`
		Description string `form:"description" binding:"max=500"`
		CategoryID  uint   `form:"category_id" binding:"required"`
		ImageURL    string // akan di-set dari file upload atau existing
	}
	// Bind form data
	if err := c.ShouldBind(&input); err != nil {
		helpers.BadRequest(c, "Invalid form data: "+err.Error())
		return
	}

	// 5. Validasi dan konversi category_id
	categoryID := input.CategoryID
	fmt.Printf("Input received: %+v\n", categoryID)
	if err != nil {
		helpers.BadRequest(c, "Invalid category_id format")
		return
	}

	// 6. Validasi CategoryID exists
	var categoryExists bool
	if err := s.db.Model(&models.Category{}).Select("count(*) > 0").Where("id = ?", uint(categoryID)).Find(&categoryExists).Error; err != nil {
		helpers.InternalServerError(c, "Database error")
		return
	}
	if !categoryExists {
		helpers.BadRequest(c, "Category not found")
		return
	}

	// 7. Handle file upload jika ada
	imageURL := gallery.ImageURL // default gunakan URL yang sudah ada
	
	file, header, err := c.Request.FormFile("image")
	if err == nil && header != nil {
		defer file.Close()
		
		// Validasi file type
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/jpg":  true,
			"image/png":  true,
			"image/gif":  true,
		}
		
		contentType := header.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			helpers.BadRequest(c, "Invalid file type. Only JPEG, PNG, GIF allowed")
			return
		}
		
		// Validasi file size (max 5MB)
		if header.Size > 5*1024*1024 {
			helpers.BadRequest(c, "File too large. Maximum size is 5MB")
			return
		}
		
		// Upload file (contoh ke local storage atau cloud)
		// uploadedURL, err := s.uploadImage(file, header)
		uploadedURL, err := SaveUploadedFile(file, header)
		if err != nil {
			helpers.InternalServerError(c, "Failed to upload image: "+err.Error())
			return
		}
		
		imageURL = uploadedURL
	}

	// 8. Update gallery dengan error handling
	updateData := models.Gallery{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    imageURL,
		CategoryID:  categoryID,
	}

	if err := s.db.Model(&gallery).Updates(updateData).Error; err != nil {
		helpers.InternalServerError(c, "Failed to update gallery")
		return
	}

	// 9. Reload data yang sudah diupdate untuk response
	if err := s.db.First(&gallery, gallery.ID).Error; err != nil {
		helpers.InternalServerError(c, "Failed to reload gallery data")
		return
	}

	// Broadcast update
	BroadcastUpdateGallery(map[string]interface{}{
		"ID":          gallery.ID,
		"title":       gallery.Title,
		"description": gallery.Description,
		"image_url":   gallery.ImageURL,
		"category_id": gallery.CategoryID,
		"updated_at":  gallery.UpdatedAt,
	})

	helpers.Success(c, "Gallery updated successfully", gallery)
}

// Helper function untuk upload image
/*
func (s *Server) uploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), ext)
	
	// Create uploads directory if not exists
	uploadDir := "./uploads/galleries"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}
	
	// Save file to local storage
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	
	// Return URL (adjust based on your static file serving setup)
	return fmt.Sprintf("/uploads/galleries/%s", filename), nil
}
*/

// Helper function untuk generate random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
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
	// Broadcast delete
	BroadcastDeleteGallery(id)

	helpers.Success(c, "Gallery deleted", gallery)	
}
