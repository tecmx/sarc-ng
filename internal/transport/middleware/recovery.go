package middleware

import (
	"github.com/gin-gonic/gin"
)

// Recovery returns a recovery middleware that handles panics
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
