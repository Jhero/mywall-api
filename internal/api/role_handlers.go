package api

import (
	"mywall-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mywall-api/internal/helpers"
)

type RoleRequest struct {
	ID    			string `json:"ID" binding:"required,max=20"`
	Name 			string `json:"name" binding:"required,max=50"`
	Description 	string `json:"description" binding:"required,max=200"`
}

func (s *Server) getRoles(c *gin.Context) {
	userID := c.GetUint("user_id")
	var roles []models.Role
	s.db.Where("user_id = ?", userID).Find(&roles)
	helpers.Success(c, "Roles retrieved successfully", roles)	
}

func (s *Server) getRole(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var role models.Role
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&role).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	helpers.Success(c, "Role retrieved successfully", role)	
}

func (s *Server) createRole(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				switch e.Field() {
					case "Name":
						if e.Tag() == "required" {
							errorMessages["name"] = "Name is required" 
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
	
	role := models.Role{
		ID:       		req.ID,
		Name: 			req.Name,
		UserID:      	userID,
		// Set other fields as needed
	}
	if result := s.db.Create(&role); result.Error != nil {
		helpers.InternalServerError(c,"Failed to create role")
		return
	}
	helpers.Created(c,"Role created successfully",role)
}

func (s *Server) updateRole(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var role models.Role
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&role).Error; err != nil {
		helpers.NotFound(c,"Role not found")
		return
	}

	var input models.Role
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.BadRequest(c,err.Error())
		return
	}

	s.db.Model(&role).Updates(models.Role{
		ID:       	input.ID,
		Name: 		input.Name,
	})
	helpers.Success(c, "Role updated successfully", role)
}

func (s *Server) deleteRole(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	var role models.Role
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&role).Error; err != nil {
		helpers.NotFound(c, "Role not found")
		return
	}
	s.db.Delete(&role)
	helpers.Success(c, "Role deleted", role)	
}
