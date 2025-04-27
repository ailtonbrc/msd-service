// internal/api/middleware/permission_middleware.go
package middleware

import (
	"net/http"

	"clinica_server/internal/auth"

	"github.com/gin-gonic/gin"
)

// RequirePermission cria um middleware que verifica se o usuário tem a permissão necessária
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter o usuário do contexto (já definido pelo middleware de autenticação)
		userClaims, exists := auth.GetUserFromContext(c.Request.Context())
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Não autorizado",
				"message": "Usuário não autenticado",
			})
			c.Abort()
			return
		}

		// Verificar se o usuário tem a permissão necessária
		if !auth.HasPermission(userClaims, permission) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Acesso proibido",
				"message": "Você não tem permissão para acessar este recurso",
			})
			c.Abort()
			return
		}

		// Continuar para o próximo handler
		c.Next()
	}
}

// RequireRole cria um middleware que verifica se o usuário tem o papel necessário
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter o usuário do contexto (já definido pelo middleware de autenticação)
		userClaims, exists := auth.GetUserFromContext(c.Request.Context())
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Não autorizado",
				"message": "Usuário não autenticado",
			})
			c.Abort()
			return
		}

		// Verificar se o usuário tem o papel necessário
		if !auth.HasRole(userClaims, role) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Acesso proibido",
				"message": "Você não tem o papel necessário para acessar este recurso",
			})
			c.Abort()
			return
		}

		// Continuar para o próximo handler
		c.Next()
	}
}

// RequireScope cria um middleware que verifica se o usuário tem o escopo necessário
func RequireScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter o usuário do contexto (já definido pelo middleware de autenticação)
		userClaims, exists := auth.GetUserFromContext(c.Request.Context())
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Não autorizado",
				"message": "Usuário não autenticado",
			})
			c.Abort()
			return
		}

		// Verificar se o usuário tem o escopo necessário
		if !auth.HasScope(userClaims, scope) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Acesso proibido",
				"message": "Você não tem o escopo necessário para acessar este recurso",
			})
			c.Abort()
			return
		}

		// Continuar para o próximo handler
		c.Next()
	}
}