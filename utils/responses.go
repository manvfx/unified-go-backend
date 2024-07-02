package utils

import "github.com/gin-gonic/gin"

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// PaginatedResponse represents the structure of a paginated response.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// CreateErrorResponse creates an error response.
func CreateErrorResponse(message string, validationErrors map[string]string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Errors:  validationErrors,
	}
}

// CreatePaginatedResponse creates a paginated response.
func CreatePaginatedResponse(data interface{}, page, limit, totalCount int) PaginatedResponse {
	totalPages := totalCount / limit
	if totalCount%limit > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}

// JSONErrorResponse writes a JSON error response.
func JSONErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, CreateErrorResponse(message, nil))
}
