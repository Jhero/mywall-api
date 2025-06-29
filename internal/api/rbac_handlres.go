package api

import (
	"mywall-api/internal/models"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mywall-api/internal/helpers"
	"fmt"
)

type RbacRequest struct {
	MenuID    	string 				`json:"menu_id" binding:"required,max=50"`
	Permission 	PermissionStruct 	`json:"permission" binding:"required"`
	RoleID 		string 				`json:"role_id" binding:"required,max=20"`
}

// Permission structure
type PermissionStruct struct {
	Read    bool   `json:"read"`
	Edit    bool   `json:"edit"`
	Delete  bool   `json:"delete"`
	Create  bool   `json:"create"`
	Search  bool   `json:"search"`
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
					case "RoleID":
						if e.Tag() == "required" {
							errorMessages["role_id"] = "Role ID is required" 
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

	// Check if role exists for role_id
	var role models.Role
	roleID := req.RoleID
	fmt.Println("Role ID:", roleID)
	if result := s.db.Where("id = ?", roleID).First(&role); result.Error != nil {
		helpers.NotFound(c,"Invalid role")
		return
	}

	// Convert permission struct to JSON string
	permissionJSON, err := json.Marshal(req.Permission)
	if err != nil {
		helpers.InternalServerError(c, "Failed to process permission data")
		return
	}
	
	rbac := models.Rbac{
		MenuID:       	req.MenuID,
		Permission: 	string(permissionJSON),
		UserID:      	userID,
		RoleID:      	req.RoleID,
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
