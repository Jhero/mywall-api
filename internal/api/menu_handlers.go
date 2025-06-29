package api

import (
	"mywall-api/internal/models"
	"net/http"
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mywall-api/internal/helpers"
)

type MenuRequest struct {
	ID    	string `json:"ID" binding:"required,max=50"`
	Path 	string `json:"path" binding:"required,max=150"`
}

func (s *Server) getMenus(c *gin.Context) {
	userID := c.GetUint("user_id")
	// fmt.Println("User ID:", userID)
	var menus []models.Menu
	s.db.Where("user_id = ?", userID).Find(&menus)
	helpers.Success(c, "Menus retrieved successfully", menus)	
}

func (s *Server) getMenu(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var menu models.Menu
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}
	helpers.Success(c, "Menu retrieved successfully", menu)	
}

func (s *Server) createMenu(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req MenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				switch e.Field() {
					case "Path":
						if e.Tag() == "required" {
							errorMessages["path"] = "Path is required" 
						}
					case "ID":
						if e.Tag() == "required" {
							errorMessages["id"] = "Id is required" 
						}
					case "UserID":
						if e.Tag() == "required" {
							errorMessages["user_id"] = "User ID is required" 
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
		helpers.NotFound(c,"Invalid user")
		return
	}
	
	menu := models.Menu{
		ID:       		req.ID,
		Path: 			req.Path,
		UserID:      	userID,
		// Set other fields as needed
	}
	if result := s.db.Create(&menu); result.Error != nil {
		helpers.InternalServerError(c,"Failed to create menu")
		return
	}
	helpers.Created(c,"Menu created successfully",menu)
}

func (s *Server) updateMenu(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var menu models.Menu
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&menu).Error; err != nil {
		helpers.NotFound(c,"Menu not found")
		return
	}

	var input models.Menu
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c,err.Error())
		return
	}

	s.db.Model(&menu).Updates(models.Menu{
		ID:       	input.ID,
		Path: 		input.Path,
	})
	helpers.Success(c, "Menu updated successfully", menu)
}

func (s *Server) deleteMenu(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var menu models.Menu
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&menu).Error; err != nil {
		helpers.NotFound(c, "Menu not found")
		return
	}
	s.db.Delete(&menu)
	helpers.Success(c, "Menu deleted", menu)	
}
