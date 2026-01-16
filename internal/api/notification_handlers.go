package api

import (
    "encoding/json"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "mywall-api/internal/models"
    "mywall-api/internal/helpers"
)
	
// PERBAIKAN: Struct NotificationHandlers hanya berisi db
type NotificationHandlers struct {
    db *gorm.DB
}

// listNotifications: ambil semua notifikasi user
func (h *Server) listNotifications(c *gin.Context) {
    userIDVal, exists := c.Get("user_id")
    if !exists {	
        helpers.Unauthorized(c, "unauthorized")
        return
    }
    userID := userIDVal.(uint)

    var notifs []models.Notification
    if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifs).Error; err != nil {
        helpers.InternalServerError(c, "failed to fetch notifications")
        return
    }
    helpers.Success(c, "notifications fetched", notifs)
}

// PERBAIKAN: Method CreateNotificationDirect sekarang bisa akses h.db
func (h *NotificationHandlers) CreateNotificationDirect(userID uint, title, body, notifType string, metadata map[string]interface{}) error {
    metadataJSON, err := json.Marshal(metadata)
    if err != nil {
        return err
    }

    n := models.Notification{
        ID:       uuid.New().String(),
        UserID:   userID,
        Title:    title,
        Body:     body,
        Type:     notifType,
        Metadata: string(metadataJSON),
        IsRead:   false,
    }
    
    if err := h.db.Create(&n).Error; err != nil {
        return err
    }

    BroadcastNotification(map[string]interface{}{
        "id":       n.ID,
        "user_id":  n.UserID,
        "title":    n.Title,
        "body":     n.Body,
        "type":     n.Type,
        "metadata": n.Metadata,
        "is_read":  n.IsRead,
    })
    
    log.Printf("Broadcasted new notification for user %d", userID)
    return nil
}

// createNotification: buat notifikasi baru dan broadcast via WebSocket
func (h *Server) createNotification(c *gin.Context) {
    var input struct {
        UserID   uint                   `json:"userId"`  // PERBAIKAN: Ubah ke uint langsung
        Title    string                 `json:"title"`
        Body     string                 `json:"body"`
        Type     string                 `json:"type"`
        Metadata map[string]interface{} `json:"metadata"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        helpers.BadRequest(c, "invalid payload")
        return
    }

    // PERBAIKAN: Panggil CreateNotificationDirect dengan benar
    notifHandler := &NotificationHandlers{db: h.db}
    if err := notifHandler.CreateNotificationDirect(input.UserID, input.Title, input.Body, input.Type, input.Metadata); err != nil {
        helpers.InternalServerError(c, "failed to create notification")
        return
    }

    helpers.Success(c, "Notification created successfully", map[string]interface{}{
        "user_id":  input.UserID,
        "title":    input.Title,
        "body":     input.Body,
        "type":     input.Type,
        "metadata": input.Metadata,
        "is_read":  false,
    })
}

// markRead: tandai notifikasi sebagai dibaca
func (h *Server) markRead(c *gin.Context) {
    var input struct {
        UserID  uint   `json:"userId"`   // PERBAIKAN: Ubah ke uint
        NotifID string `json:"notifId"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        helpers.BadRequest(c, "invalid payload")
        return
    }

    if err := h.db.Model(&models.Notification{}).
        Where("id = ? AND user_id = ?", input.NotifID, input.UserID).
        Update("is_read", true).Error; err != nil {
        helpers.InternalServerError(c, "failed to mark read")
        return
    }

    // Broadcast badge update
    var unreadCount int64
    h.db.Model(&models.Notification{}).
        Where("user_id = ? AND is_read = false", input.UserID).
        Count(&unreadCount)
    BroadcastNotification(map[string]interface{}{
        "user_id": input.UserID,
        "unread":  unreadCount,
    })
    helpers.Success(c, "notification marked read", map[string]interface{}{
        "user_id": input.UserID,
        "unread":  unreadCount,
    })
}