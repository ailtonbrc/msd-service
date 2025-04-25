package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Erros de autenticação
var (
	ErrAuthHeaderMissing = errors.New("authorization header is required")
	ErrInvalidAuthHeader = errors.New("invalid authorization header format")
	ErrInvalidToken      = errors.New("invalid or expired token")
)

// Claims representa as claims do JWT
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware gerencia a autenticação via JWT
type AuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
}

// JWTAuthMiddleware implementa AuthMiddleware usando JWT
type JWTAuthMiddleware struct {
	secretKey []byte
}

// NewJWTAuthMiddleware cria uma nova instância de JWTAuthMiddleware
func NewJWTAuthMiddleware(secretKey string) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		secretKey: []byte(secretKey),
	}
}

// RequireAuth é um middleware que verifica se o usuário está autenticado
func (m *JWTAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrAuthHeaderMissing.Error()})
			return
		}

		// Verificar formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidAuthHeader.Error()})
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		// Parse do token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verificar algoritmo de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
			}
			return m.secretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			return
		}

		// Armazenar claims no contexto para uso posterior
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// GenerateToken gera um novo token JWT
func (m *JWTAuthMiddleware) GenerateToken(userID, role string, expirationTime time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}
