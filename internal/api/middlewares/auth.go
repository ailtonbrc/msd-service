package middlewares

import (
	"net/http"
	"strings"

	"clinica_server/config"
	"clinica_server/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Token não fornecido")
			c.Abort()
			return
		}

		// Extrair token do cabeçalho
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Formato de token inválido")
			c.Abort()
			return
		}

		// Validar token
		claims, err := utils.ValidateToken(parts[1], cfg)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", err.Error())
			c.Abort()
			return
		}

		// Armazenar claims no contexto
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("Perfil", claims.Perfil)

		c.Next()
	}
}
