package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response representa a estrutura padrão de resposta da API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// SuccessResponse envia uma resposta de sucesso
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse envia uma resposta de erro
func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// ValidationErrorResponse envia uma resposta de erro de validação
func ValidationErrorResponse(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Data:    errors,
	})
}
