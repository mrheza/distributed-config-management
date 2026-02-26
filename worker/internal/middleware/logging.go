package middleware

import (
	sharedmiddleware "github.com/mrheza/distributed-config-management/shared/middleware"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return sharedmiddleware.RequestLogger()
}
