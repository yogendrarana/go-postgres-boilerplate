package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Token     string    `gorm:"type:varchar(255);unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

// Get refresh tokens by user ID
func GetRefreshTokenByUserID(db *gorm.DB, userID uint) (*[]RefreshToken, error) {
	var refreshToken []RefreshToken
	err := db.Where("user_id = ?", userID).Find(&refreshToken).Error
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// Delete refresh token by token hash
func DeleteRefreshToken(db *gorm.DB, tokenHash string) error {
	var refreshToken RefreshToken
	err := db.Where("token_hash = ?", tokenHash).Delete(&refreshToken).Error
	if err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		// No records were deleted, handle it as needed (return an error, etc.)
		return errors.New("refresh token not found")
	}

	return nil
}
