package routes

import (
	"go-gin-postgres/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(authGroup *gin.RouterGroup) {
	authGroup.POST("/register", handlers.RegisterWithEmailPassword)
}
