package middlewares

import (
	db "go-gin-postgres/internal/db"

	"github.com/gin-gonic/gin"
)

func DBMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		database := db.GetDB()
		ctx.Set("db", database)
		ctx.Next()
	}
}
