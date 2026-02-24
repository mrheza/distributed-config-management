package middleware

import (
	"controller/internal/httpresponse"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != key {
			httpresponse.Unauthorized(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
