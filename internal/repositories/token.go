package repositories

import (
	"saaa-api/internal/models"
	"time"

	"gorm.io/gorm"
)

// RefreshTokenRepository refresh token repository interface
type RefreshTokenRepository interface {
	Create(database *gorm.DB, i interface{}) error
	Update(database *gorm.DB, i interface{}) error
	Delete(database *gorm.DB, i interface{}) error
	FindByID(database *gorm.DB, id uint, i interface{}) error
	FindByRefreshTokenExpired(database *gorm.DB, refreshTokenString string, refreshToken *models.RefreshToken) error
	FindByRefreshTokenWithoutCondition(database *gorm.DB, refreshTokenString string, refreshToken *models.RefreshToken) error
	Remove(database *gorm.DB, refreshToken string) error
}

type refreshTokenRepository struct {
	Repository
}

// NewRefreshTokenRepository for new refresh token repository.
func NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepository{
		NewRepository(),
	}
}

// FindByRefreshTokenExpired find refreshToken by refreshToken expired
func (repo *refreshTokenRepository) FindByRefreshTokenExpired(database *gorm.DB, refreshTokenString string, refreshToken *models.RefreshToken) error {
	if err := database.
		Where("refresh_token = ? AND used = ? AND expire_time <= ?", refreshTokenString, false, time.Now()).
		Last(refreshToken).Error; err != nil {
		return err
	}

	return nil
}

// FindByRefreshTokenWithoutCondition find refreshToken by refreshToken without condition
func (repo *refreshTokenRepository) FindByRefreshTokenWithoutCondition(database *gorm.DB, refreshTokenString string, refreshToken *models.RefreshToken) error {
	if err := database.
		Where("refresh_token = ?", refreshTokenString).
		First(refreshToken).Error; err != nil {
		return err
	}

	return nil
}

// Remove update field used equal true
func (repo *refreshTokenRepository) Remove(database *gorm.DB, refreshTokenString string) error {
	if err := database.
		Model(&models.RefreshToken{}).
		Where("refresh_token = ?", refreshTokenString).
		Update("used", true).Error; err != nil {
		return err
	}

	return nil
}
