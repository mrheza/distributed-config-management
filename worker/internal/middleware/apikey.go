package middleware

import (
	sharedmiddleware "github.com/mrheza/distributed-config-management/shared/middleware"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(key string) gin.HandlerFunc {
	return sharedmiddleware.APIKeyAuth(key)
}
