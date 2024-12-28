package handlers

import (
	"go-gin-postgres/internal/database/models"
	"go-gin-postgres/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const bcryptCost = 10

type RegisterInput struct {
	FullName        string `json:"full_name" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

// Helper function to send JSON error responses
func sendErrorResponse(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"success": false, "message": message})
}

// RegisterWithEmailPassword handles user registration
func RegisterWithEmailPassword(ctx *gin.Context) {
	var input RegisterInput

	// Get DB connection
	db := ctx.MustGet("db").(*gorm.DB)

	// Bind JSON input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		sendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Validate password confirmation
	if input.Password != input.ConfirmPassword {
		sendErrorResponse(ctx, http.StatusBadRequest, "Passwords do not match")
		return
	}

	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if user already exists
	var existingUser models.User
	if err := tx.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		sendErrorResponse(ctx, http.StatusBadRequest, "Email already registered.")
		return
	}

	// Create new user
	user := models.User{
		FullName: input.FullName,
		Email:    input.Email,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		sendErrorResponse(ctx, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Hash password and store it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		tx.Rollback()
		sendErrorResponse(ctx, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	password := models.Password{
		UserID: user.ID,
		Hash:   string(hashedPassword),
	}
	if err := tx.Create(&password).Error; err != nil {
		tx.Rollback()
		sendErrorResponse(ctx, http.StatusInternalServerError, "Failed to save password")
		return
	}

	// Generate access token
	signedAccessToken, err := services.GenerateAccessToken(user.ID)
	if err != nil {
		tx.Rollback()
		sendErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	// Generate refresh token
	refreshToken, hashedToken, err := services.GenerateRefreshTokenAndHash()
	if err != nil {
		tx.Rollback()
		sendErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Save refresh token in the database
	dbRefreshToken := models.RefreshToken{
		UserID: user.ID,
		Token:  hashedToken,
	}
	if err := tx.Create(&dbRefreshToken).Error; err != nil {
		tx.Rollback()
		sendErrorResponse(ctx, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		sendErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	// Set HTTP-only cookie for refresh token
	ctx.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data": gin.H{
			"access_token": signedAccessToken,
			"user":         user,
		},
	})
}
