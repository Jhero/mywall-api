package auth

import (
	"crypto/rand"
	"encoding/hex"
	"mywall-api/internal/models"
	"gorm.io/gorm"
)

type APIKeyService struct {
	keyLength int
}

func NewAPIKeyService(keyLength int) *APIKeyService {
	return &APIKeyService{
		keyLength: keyLength,
	}
}

func (a *APIKeyService) GenerateAPIKey() (string, error) {
	bytes := make([]byte, a.keyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *APIKeyService) ValidateAPIKey(db *gorm.DB, apiKey string) (*models.User, error) {
	var user models.User
	if err := db.Where("api_key = ? AND is_active = ?", apiKey, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
