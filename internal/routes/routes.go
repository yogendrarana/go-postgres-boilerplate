package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup) {
	AuthRoutes(router.Group("/auth"))
}
