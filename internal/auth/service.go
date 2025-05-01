package auth

import (
	"errors"
	"mywall-api/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service handles authentication-related operations
type Service struct {
	db          *gorm.DB
	jwtService  *JWTService
	apiKeyService *APIKeyService
}

// NewService creates a new auth service
func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{
		db:          db,
		jwtService:  NewJWTService(jwtSecret),
		apiKeyService: NewAPIKeyService(32), // 32 bytes for API key
	}
}

func (s *Service) Register(email, password, name string) (*models.User, error) {
	// Check if user exists
	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate API key
	apiKey, err := s.apiKeyService.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
		ApiKey:   apiKey,
		Role:     "user",
		IsActive: true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(email, password string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(&user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) ValidateJWT(token string) (*models.User, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) ValidateAPIKey(apiKey string) (*models.User, error) {
	return s.apiKeyService.ValidateAPIKey(s.db, apiKey)
}

func (s *Service) RegenerateAPIKey(userID uint) (string, error) {
	apiKey, err := s.apiKeyService.GenerateAPIKey()
	if err != nil {
		return "", err
	}

	if err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("api_key", apiKey).Error; err != nil {
		return "", err
	}

	return apiKey, nil
}
