package middlewares

import (
	"go-gin-postgres/internal/database"

	"github.com/gin-gonic/gin"
)

func DBMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		database := database.GetDB()
		ctx.Set("db", database)
		ctx.Next()
	}
}
