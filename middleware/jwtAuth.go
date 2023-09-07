package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mulukenhailu/Diary_api/helper"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := helper.ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Restricted"})
			context.Abort()
			return
		}
		context.Next()
	}
}
