package handlers

import (
	"net/http"

	"go-gin-postgres/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Helper function to send JSON error responses
func sendErrorResponse(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{"success": false, "message": message})
}

// RegisterWithEmailPassword handles user registration
func RegisterWithEmailPassword(ctx *gin.Context) {
	var input services.RegisterInput

	// Get DB connection
	db := ctx.MustGet("db").(*gorm.DB)

	// Bind JSON input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		sendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, accessToken, refreshToken, err := services.RegisterWithEmailPassword(db, input)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrEmailAlreadyRegistered {
			status = http.StatusBadRequest
		}
		sendErrorResponse(ctx, status, err.Error())
		return
	}

	// Set HTTP-only cookie for refresh token
	ctx.SetCookie("refresh_token", refreshToken, 24*60*60, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data": gin.H{
			"access_token": accessToken,
			"user":         user,
		},
	})
}
