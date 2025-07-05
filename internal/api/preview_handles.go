package api

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/gin-gonic/gin"
)

func (s *Server) serveImage(c *gin.Context) {
	// fmt.Println("Masuk-1")
    year := c.Param("year")
    month := c.Param("month")
    day := c.Param("day")
    filename := c.Param("filename")
    
    // Basic validation
    if year == "" || month == "" || day == "" || filename == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
        return
    }
    
    // Security: prevent directory traversal attacks
    if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
        return
    }
    
    // Construct the file path
    imagePath := filepath.Join("uploads", year, month, day, filename)
    fmt.Println("path",imagePath)
    // Check if file exists
    fileInfo, err := os.Stat(imagePath)
    if os.IsNotExist(err) {
        c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
        return
    }
    
    // Optional: Set appropriate headers
    c.Header("Content-Type", "image/jpeg") // or detect MIME type
    c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
    c.Header("Cache-Control", "public, max-age=31536000") // Cache for 1 year
    
    // Serve the file
    c.File(imagePath)
}