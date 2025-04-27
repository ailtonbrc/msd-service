// internal/api/middleware/auth_middleware.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"clinica_server/internal/auth"
)

// GetJWTService retorna o serviço JWT usado pelo middleware
func (m *AuthMiddleware) GetJWTService() auth.JWTService {
    return m.jwtService
}

// AuthMiddleware é responsável por verificar a autenticação do usuário
type AuthMiddleware struct {
	jwtService auth.JWTService
	logger     *zap.Logger
}

// NewAuthMiddleware cria uma nova instância de AuthMiddleware
func NewAuthMiddleware(jwtService auth.JWTService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// RequireAuth cria um middleware que verifica se o usuário está autenticado
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair o token do cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Não autorizado",
				"message": "Token de autenticação não fornecido",
			})
			c.Abort()
			return
		}

		// Verificar o formato do token (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Não autorizado",
				"message": "Formato de token inválido",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar o token
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			m.logger.Warn("Token inválido",
				zap.Error(err),
				zap.String("token", tokenString),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Não autorizado",
				"message": "Token inválido ou expirado",
			})
			c.Abort()
			return
		}

		// Verificar se o token está na lista negra (tokens revogados)
		isRevoked, err := m.jwtService.IsTokenRevoked(tokenString)
		if err != nil {
			m.logger.Error("Erro ao verificar revogação de token",
				zap.Error(err),
				zap.String("token", tokenString),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Erro interno",
				"message": "Erro ao verificar autenticação",
			})
			c.Abort()
			return
		}

		if isRevoked {
			m.logger.Warn("Token revogado",
				zap.String("token", tokenString),
				zap.Uint("user_id", claims.UserID),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Não autorizado",
				"message": "Token revogado",
			})
			c.Abort()
			return
		}

		// Adicionar as informações do usuário ao contexto
		ctx := context.WithValue(c.Request.Context(), auth.UserClaimsKey, claims)
		c.Request = c.Request.WithContext(ctx)

		// Registrar o acesso
		m.logger.Info("Usuário autenticado",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		// Continuar para o próximo handler
		c.Next()
	}
}

// OptionalAuth cria um middleware que verifica a autenticação, mas não exige
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair o token do cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Sem token, continuar sem autenticação
			c.Next()
			return
		}

		// Verificar o formato do token (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Token mal formatado, continuar sem autenticação
			c.Next()
			return
		}

		tokenString := parts[1]

		// Validar o token
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			// Token inválido, continuar sem autenticação
			c.Next()
			return
		}

		// Verificar se o token está na lista negra (tokens revogados)
		isRevoked, err := m.jwtService.IsTokenRevoked(tokenString)
		if err != nil || isRevoked {
			// Erro ou token revogado, continuar sem autenticação
			c.Next()
			return
		}

		// Adicionar as informações do usuário ao contexto
		ctx := context.WithValue(c.Request.Context(), auth.UserClaimsKey, claims)
		c.Request = c.Request.WithContext(ctx)

		// Registrar o acesso
		m.logger.Info("Usuário autenticado (opcional)",
			zap.Uint("user_id", claims.UserID),
			zap.String("username", claims.Username),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		// Continuar para o próximo handler
		c.Next()
	}
}