package utils

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func CreateErrorResponse(message string) gin.H {
	return gin.H{"error": message}
}
