package api

import (
	"mywall-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"	
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				// Generate specific error messages based on field and validation tag
				switch e.Field() {
					case "Email":
						if e.Tag() == "required" {
							errorMessages["email"] = "Email is required"
						} else if e.Tag() == "email" {
							errorMessages["email"] = "Invalid email format"
						}
					case "Password":
						if e.Tag() == "required" {
							errorMessages["password"] = "Password is required"
						} else if e.Tag() == "min" {
							errorMessages["password"] = "Password must be at least 8 characters long"
						}
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

	user, err := s.auth.Register(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate initial JWT token
	token, err := s.auth.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"api_key": user.ApiKey,
		},
	})
}

func (s *Server) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	token, err := s.auth.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Get user info
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"api_key": user.ApiKey,
		},
	})
}

func (s *Server) handleRegenerateApiKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	apiKey, err := s.auth.RegenerateAPIKey(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to regenerate API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": apiKey})
}
