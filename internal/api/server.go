package api

import (
	"mywall-api/internal/auth"
	// "mywall/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mywall-api/internal/helpers"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	db     *gorm.DB
	auth   *auth.Service
}

// NewServer creates a new server instance
func NewServer(db *gorm.DB, auth *auth.Service) *Server {
	server := &Server{
		router: gin.Default(),
		db:     db,
		auth:   auth,
	}
	server.setupRoutes()
	return server
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Public auth routes
	authRoutes := s.router.Group("/auth")
	{
		authRoutes.POST("/register", s.handleRegister)
		authRoutes.POST("/login", s.handleLogin)
	}

	// Protected routes
	apiRoutes := s.router.Group("/api")
	apiRoutes.Use(s.authMiddleware())
	{
		// API key management
		apiRoutes.POST("/regenerate-api-key", s.handleRegenerateApiKey)

		apiRoutes.GET("/images/:year/:month/:day/:filename", s.serveImage)

		// Other API routes
		apiRoutes.GET("/galleries", s.getGalleries)
		apiRoutes.POST("/galleries", s.createGallery)
		apiRoutes.GET("/galleries/:id", s.getGallery)
		apiRoutes.PUT("/galleries/:id", s.updateGallery)
		apiRoutes.DELETE("/galleries/:id", s.deleteGallery)

		apiRoutes.GET("/categories", s.getCategories)
		apiRoutes.POST("/categories", s.createCategory)
		apiRoutes.GET("/categories/:id", s.getCategory)
		apiRoutes.PUT("/categories/:id", s.updateCategory)
		apiRoutes.DELETE("/categories/:id", s.deleteCategory)

		apiRoutes.GET("/menus", s.getMenus)
		apiRoutes.POST("/menus", s.createMenu)
		apiRoutes.GET("/menus/:id", s.getMenu)
		apiRoutes.PUT("/menus/:id", s.updateMenu)
		apiRoutes.DELETE("/menus/:id", s.deleteMenu)


		apiRoutes.GET("/rbacs", s.getRbacs)
		apiRoutes.POST("/rbacs", s.createRbac)
		apiRoutes.GET("/rbacs/:id", s.getRbac)
		apiRoutes.PUT("/rbacs/:id", s.updateRbac)
		apiRoutes.DELETE("/rbacs/:id", s.deleteRbac)

		apiRoutes.GET("/roles", s.getRoles)
		apiRoutes.POST("/roles", s.createRole)
		apiRoutes.GET("/roles/:id", s.getRole)
		apiRoutes.PUT("/roles/:id", s.updateRole)
		apiRoutes.DELETE("/roles/:id", s.deleteRole)

	}
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	return s.router.Run(":" + port)
}

// Authentication middleware
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check for JWT token in Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				user, err := s.auth.ValidateJWT(token)
				if err == nil {
					c.Set("user_id", user.ID)
					c.Next()
					return
				}
			}
		}

		// If no valid JWT, check for API key
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			user, err := s.auth.ValidateAPIKey(apiKey)
			if err == nil {
				c.Set("user_id", user.ID)
				c.Next()
				return
			}
		}
		helpers.Unauthorized(c, "Invalid authorization token or API key required")
		c.Abort()
	}
}
