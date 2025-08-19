package repository

import (
	"go-gin-postgres/internal/models"

	"gorm.io/gorm"
)

func FindUserByEmail(tx *gorm.DB, email string) (*models.User, error) {
	var user models.User
	err := tx.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func CreatePassword(tx *gorm.DB, password *models.Password) error {
	return tx.Create(password).Error
}

func CreateRefreshToken(tx *gorm.DB, token *models.RefreshToken) error {
	return tx.Create(token).Error
}
