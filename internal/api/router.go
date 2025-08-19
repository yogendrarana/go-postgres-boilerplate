package api

import (
	middlewares "go-gin-postgres/internal/middleware"
	"go-gin-postgres/internal/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	// Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.DBMiddleware())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// TODO: add rate limiter middleware
	// TODO: serve static file

	// register all routes
	api := router.Group("/api")
	RegisterV1Routes(api)

	// home page
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return router
}

func RegisterV1Routes(router *gin.RouterGroup) {
	routes.RegisterAuthRoutes(router.Group("/auth"))
}
