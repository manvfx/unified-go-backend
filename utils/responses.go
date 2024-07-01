package utils

import (
	"github.com/gin-gonic/gin"
)

func CreateErrorResponse(message string) gin.H {
	return gin.H{"error": message}
}
