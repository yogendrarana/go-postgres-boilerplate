package services

import (
	"errors"
	"time"

	"go-gin-postgres/internal/models"
	"go-gin-postgres/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrEmailAlreadyRegistered = errors.New("Email already registered.")

type RegisterInput struct {
	FullName        string `json:"full_name" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

func RegisterWithEmailPassword(db *gorm.DB, input RegisterInput) (models.User, string, string, error) {
	var emptyUser models.User

	if input.Password != input.ConfirmPassword {
		return emptyUser, "", "", errors.New("Passwords do not match")
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if user already exists
	if _, err := repository.FindUserByEmail(tx, input.Email); err == nil {
		tx.Rollback()
		return emptyUser, "", "", ErrEmailAlreadyRegistered
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return emptyUser, "", "", err
	}

	// Create user
	user := models.User{FullName: input.FullName, Email: input.Email}
	if err := repository.CreateUser(tx, &user); err != nil {
		tx.Rollback()
		return emptyUser, "", "", err
	}

	// Store password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		tx.Rollback()
		return emptyUser, "", "", err
	}
	pwd := models.Password{UserID: user.ID, Hash: string(hash)}
	if err := repository.CreatePassword(tx, &pwd); err != nil {
		tx.Rollback()
		return emptyUser, "", "", err
	}

	// Tokens
	accessToken, err := GenerateAccessToken(user.ID)
	if err != nil {
		tx.Rollback()
		return emptyUser, "", "", err
	}
	refreshToken, hashedToken, err := GenerateRefreshTokenAndHash()
	if err != nil {
		tx.Rollback()
		return emptyUser, "", "", err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	dbRefresh := models.RefreshToken{UserID: user.ID, Token: hashedToken, ExpiresAt: expiresAt}
	if err := repository.CreateRefreshToken(tx, &dbRefresh); err != nil {
		tx.Rollback()
		return emptyUser, "", "", err
	}

	if err := tx.Commit().Error; err != nil {
		return emptyUser, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
