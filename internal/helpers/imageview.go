package helpers

import (
	"time"
	
	"gorm.io/gorm"
)

type ImageView struct {
	ID        uint      `gorm:"primaryKey"`
	GalleryID uint      `gorm:"not null"`
	Count     int       `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Server atau struct yang mengandung db
type ImageViewService struct {
	db *gorm.DB
}

func NewImageViewService(db *gorm.DB) *ImageViewService {
	return &ImageViewService{db: db}
}

// CreateOrUpdateImageView membuat atau mengupdate image view
func (s *ImageViewService) CreateOrUpdateImageView(galleryID uint, count int) (*ImageView, error) {
	var imageView ImageView
	
	// Mulai transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// Cari existing record
	result := tx.Where("gallery_id = ?", galleryID).First(&imageView)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Buat baru jika tidak ditemukan
			imageView = ImageView{
				GalleryID: galleryID,
				Count:     count,
			}
			if err := tx.Create(&imageView).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// Error lainnya
			tx.Rollback()
			return nil, result.Error
		}
	} else {
		// Update existing
		imageView.Count += count
		if err := tx.Save(&imageView).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	
	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	
	return &imageView, nil
}

// GetImageViewByGalleryID mengambil image view berdasarkan gallery ID
func (s *ImageViewService) GetImageViewByGalleryID(galleryID uint) (*ImageView, error) {
	var imageView ImageView
	result := s.db.Where("gallery_id = ?", galleryID).First(&imageView)
	if result.Error != nil {
		return nil, result.Error
	}
	return &imageView, nil
}

// UpdateImageViewCount mengupdate count image view
func (s *ImageViewService) UpdateImageViewCount(galleryID uint, additionalCount int) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	var imageView ImageView
	if err := tx.Where("gallery_id = ?", galleryID).First(&imageView).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	imageView.Count += additionalCount
	if err := tx.Save(&imageView).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}