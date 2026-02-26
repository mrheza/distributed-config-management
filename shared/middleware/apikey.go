package middleware

import (
	sharedhttpresponse "github.com/mrheza/distributed-config-management/shared/httpresponse"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != key {
			sharedhttpresponse.Unauthorized(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
