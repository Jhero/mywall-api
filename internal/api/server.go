package api

import (
	"log"
	"mywall-api/internal/auth"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"mywall-api/internal/helpers"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	db     *gorm.DB
	auth   *auth.Service
	ws     *WebSocketManager 
}

// WebSocketManager manages WebSocket connections
type WebSocketManager struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	mu        sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

var wsManager *WebSocketManager

func init() {
	wsManager = &WebSocketManager{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
	go wsManager.startBroadcasting()
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return wsManager
}

// NewServer creates a new server instance
func NewServer(db *gorm.DB, auth *auth.Service) *Server {
	server := &Server{
		router: gin.Default(),
		db:     db,
		auth:   auth,
		ws:     NewWebSocketManager(),
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
	// WebSocket route
	s.router.GET("/ws", s.handleWebSocket)

	// Protected routes
	apiRoutes := s.router.Group("/api")
	apiRoutes.Use(s.authMiddleware())
	{
		// API key management
		apiRoutes.POST("/regenerate-api-key", s.handleRegenerateApiKey)

		apiRoutes.GET("/images/:year/:month/:day/:filename", s.serveImage)
		
		apiRoutes.GET("/notifications", s.listNotifications)
		apiRoutes.POST("/notifications", s.createNotification)
		apiRoutes.POST("/notifications/read", s.markRead)
		
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

// WebSocket handler
func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	wsManager.mu.Lock()
	wsManager.clients[conn] = true
	wsManager.mu.Unlock()

	log.Printf("WebSocket client connected: %s", conn.RemoteAddr())
	
	// Send welcome message
	welcomeMsg := Message{
		Type:    "connected",
		Payload: "Connected to WebSocket server",
	}
	conn.WriteJSON(welcomeMsg)

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			wsManager.mu.Lock()
			delete(wsManager.clients, conn)
			wsManager.mu.Unlock()
			break
		}

		// Handle different message types from client
		switch msg.Type {
		case "join_galleries":
			log.Printf("Client joined galleries room")
		case "ping":
			// Respond to ping
			conn.WriteJSON(Message{Type: "pong", Payload: "pong"})
		}
	}
}

func (wm *WebSocketManager) startBroadcasting() {
	for {
		message := <-wm.broadcast
		wm.mu.RLock()
		for client := range wm.clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				delete(wm.clients, client)
			}
		}
		wm.mu.RUnlock()
	}
}

// Helper functions untuk broadcast
func BroadcastNewGallery(gallery map[string]interface{}) {
	message := Message{
		Type:    "new_gallery",
		Payload: gallery,
	}
	wsManager.broadcast <- message
	log.Printf("Broadcasted new gallery: %v", gallery["title"])
}

func BroadcastUpdateGallery(gallery map[string]interface{}) {
	message := Message{
		Type:    "update_gallery",
		Payload: gallery,
	}
	wsManager.broadcast <- message
	log.Printf("Broadcasted updated gallery: %v", gallery["title"])
}

func BroadcastDeleteGallery(galleryID string) {
	message := Message{
		Type:    "delete_gallery",
		Payload: map[string]string{"id": galleryID},
	}
	wsManager.broadcast <- message
	log.Printf("Broadcasted deleted gallery ID: %s", galleryID)
}

func BroadcastNotification(notification map[string]interface{}) {
	message := Message{
		Type:    "notification",
		Payload: notification,
	}
	wsManager.broadcast <- message
	log.Printf("Broadcasted notification: %v", notification["title"])
}

func BroadcastBadgeUpdate(userID uint, unreadCount int64) {
	message := Message{
		Type:    "badge_update",
		Payload: map[string]interface{}{"user_id": userID, "unread": unreadCount},
	}
	wsManager.broadcast <- message
	log.Printf("Broadcasted badge update for user %d: %d unread", userID, unreadCount)
}