package middleware

import (
	sharedmiddleware "github.com/mrheza/distributed-config-management/shared/middleware"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return sharedmiddleware.CORSMiddleware()
}
