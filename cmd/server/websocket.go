package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins for development
    },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

// Message struct untuk data yang dikirim via WebSocket
type Message struct {
    Type    string      `json:"type"`    // "new_gallery", "update_gallery", "delete_gallery"
    Payload interface{} `json:"payload"`
}

// WebSocket handler
func handleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    defer conn.Close()

    // Register client
    clients[conn] = true
    log.Printf("Client connected: %s", conn.RemoteAddr())
    
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
            delete(clients, conn)
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

// Broadcast message to all connected clients
func broadcastMessage(msgType string, payload interface{}) {
    message := Message{
        Type:    msgType,
        Payload: payload,
    }
    
    for client := range clients {
        err := client.WriteJSON(message)
        if err != nil {
            log.Printf("WebSocket write error: %v", err)
            client.Close()
            delete(clients, client)
        }
    }
}

// Function untuk broadcast new gallery
func BroadcastNewGallery(gallery map[string]interface{}) {
    broadcastMessage("new_gallery", gallery)
    log.Printf("Broadcasted new gallery: %v", gallery["title"])
}

// Function untuk broadcast updated gallery
func BroadcastUpdateGallery(gallery map[string]interface{}) {
    broadcastMessage("update_gallery", gallery)
    log.Printf("Broadcasted updated gallery: %v", gallery["title"])
}

// Function untuk broadcast deleted gallery
func BroadcastDeleteGallery(galleryID string) {
    broadcastMessage("delete_gallery", map[string]string{"id": galleryID})
    log.Printf("Broadcasted deleted gallery ID: %s", galleryID)
}