package routes

import (
	"go-gin-postgres/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(authGroup *gin.RouterGroup) {
	authGroup.POST("/auth/register", handlers.RegisterWithEmailPassword)
}
