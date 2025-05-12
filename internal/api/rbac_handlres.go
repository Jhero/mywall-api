package api

import (
	"mywall-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mywall-api/internal/helpers"
)

type RbacRequest struct {
	MenuID    	string `json:"menu_id" binding:"required,max=50"`
	Permission 	string `json:"permission" binding:"required,max=200"`
}

func (s *Server) getRbacs(c *gin.Context) {
	userID := c.GetUint("user_id")
	var rbacs []models.Rbac
	s.db.Where("user_id = ?", userID).Find(&rbacs)
	helpers.Success(c, "Rbacs retrieved successfully", rbacs)	
}

func (s *Server) getRbac(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var rbac models.Rbac
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&rbac).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rbac not found"})
		return
	}
	helpers.Success(c, "Rbac retrieved successfully", rbac)	
}

func (s *Server) createRbac(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req RbacRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				switch e.Field() {
					case "MenuID":
						if e.Tag() == "required" {
							errorMessages["menu_id"] = "Menu ID is required" 
						}
					case "Permission":
						if e.Tag() == "required" {
							errorMessages["permission"] = "Permission is required" 
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
	
	rbac := models.Rbac{
		MenuID:       	req.MenuID,
		Permission: 	req.Permission,
		UserID:      	userID,
		// Set other fields as needed
	}
	if result := s.db.Create(&rbac); result.Error != nil {
		helpers.InternalServerError(c,"Failed to create rbac")
		return
	}
	helpers.Created(c,"Rbac created successfully",rbac)
}

func (s *Server) updateRbac(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var rbac models.Rbac
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&rbac).Error; err != nil {
		helpers.NotFound(c,"Rbac not found")
		return
	}

	var input models.Rbac
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c,err.Error())
		return
	}

	s.db.Model(&rbac).Updates(models.Rbac{
		MenuID:       	input.MenuID,
		Permission: 	input.Permission,
		UserID:      	userID,
	})
	helpers.Success(c, "Rbac updated successfully", rbac)
}

func (s *Server) deleteRbac(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var rbac models.Rbac
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&rbac).Error; err != nil {
		helpers.NotFound(c, "Rbac not found")
		return
	}
	s.db.Delete(&rbac)
	helpers.Success(c, "Rbac deleted", rbac)	
}
