package helpers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// ApiResponse represents the standard API response structure
type ApiResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success sends a successful response with status code 200
func Success(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusOK, true, message, data)
}

// Created sends a successful creation response with status code 201
func Created(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusCreated, true, message, data)
}

// BadRequest sends a bad request error response with status code 400
func BadRequest(c *gin.Context, message string) {
	SendResponse(c, http.StatusBadRequest, false, message, nil)
}

// Unauthorized sends an unauthorized error response with status code 401
func Unauthorized(c *gin.Context, message string) {
	SendResponse(c, http.StatusUnauthorized, false, message, nil)
}

// Forbidden sends a forbidden error response with status code 403
func Forbidden(c *gin.Context, message string) {
	SendResponse(c, http.StatusForbidden, false, message, nil)
}

// NotFound sends a not found error response with status code 404
func NotFound(c *gin.Context, message string) {
	SendResponse(c, http.StatusNotFound, false, message, nil)
}

func NotContent(c *gin.Context, message string) {
	SendResponse(c, http.StatusNoContent, false, message, nil)
}

// InternalServerError sends a server error response with status code 500
func InternalServerError(c *gin.Context, message string) {
	SendResponse(c, http.StatusInternalServerError, false, message, nil)
}

// ValidationError sends a validation error response with status code 422
func ValidationError(c *gin.Context, message string, data interface{}) {
	SendResponse(c, http.StatusUnprocessableEntity, false, message, data)
}

// SendResponse sends a customized response with the given parameters
func SendResponse(c *gin.Context, statusCode int, success bool, message string, data interface{}) {
	c.JSON(statusCode, ApiResponse{
		Status:  success,
		Message: message,
		Data:    data,
	})
}