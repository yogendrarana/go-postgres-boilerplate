package routers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes initializes and returns a Gin router with middleware and routes
func NewRouter() *gin.Engine {
	router := gin.Default()

	// Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Gin!",
		})
	})

	return router
}
