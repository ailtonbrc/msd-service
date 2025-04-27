// internal/auth/jwt.go
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Erros relacionados a JWT
var (
	ErrInvalidToken = errors.New("token inválido")
	ErrExpiredToken = errors.New("token expirado")
	ErrTokenRevoked = errors.New("token revogado")
)

// UserClaims representa os claims personalizados para o JWT
type UserClaims struct {
	UserID    uint     `json:"user_id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles"`
	Permissions []string `json:"permissions"`
	Scopes    []string `json:"scopes"`
	jwt.RegisteredClaims
}

// JWTService define a interface para operações com JWT
type JWTService interface {
	GenerateToken(userID uint, username, email string, roles, permissions, scopes []string, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (*UserClaims, error)
	RefreshToken(tokenString string) (string, error)
	RevokeToken(tokenString string) error
	IsTokenRevoked(tokenString string) (bool, error)
}

// DefaultJWTService implementa JWTService
type DefaultJWTService struct {
	secretKey []byte
	// Aqui poderia ter um repositório para armazenar tokens revogados
}

// NewJWTService cria uma nova instância de JWTService
func NewJWTService(secretKey string) JWTService {
	return &DefaultJWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken gera um novo token JWT
func (s *DefaultJWTService) GenerateToken(
	userID uint,
	username, email string,
	roles, permissions, scopes []string,
	duration time.Duration,
) (string, error) {
	// Definir o tempo de expiração
	expirationTime := time.Now().Add(duration)

	// Criar os claims
	claims := &UserClaims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		Roles:     roles,
		Permissions: permissions,
		Scopes:    scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "clinica-tea-api",
			Subject:   username,
		},
	}

	// Criar o token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assinar o token
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valida um token JWT
func (s *DefaultJWTService) ValidateToken(tokenString string) (*UserClaims, error) {
	// Definir a função de validação de chave
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Verificar o método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	}

	// Analisar o token
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Verificar se o token é válido
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Extrair os claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken renova um token JWT
func (s *DefaultJWTService) RefreshToken(tokenString string) (string, error) {
	// Validar o token atual
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		// Permitir renovação de tokens expirados, mas não inválidos
		if !errors.Is(err, ErrExpiredToken) {
			return "", err
		}
	}

	// Verificar se o token está revogado
	isRevoked, err := s.IsTokenRevoked(tokenString)
	if err != nil {
		return "", err
	}
	if isRevoked {
		return "", ErrTokenRevoked
	}

	// Gerar um novo token com os mesmos claims, mas com nova data de expiração
	return s.GenerateToken(
		claims.UserID,
		claims.Username,
		claims.Email,
		claims.Roles,
		claims.Permissions,
		claims.Scopes,
		time.Hour*24, // Duração padrão para o novo token
	)
}

// RevokeToken revoga um token JWT
func (s *DefaultJWTService) RevokeToken(tokenString string) error {
	// Validar o token
	_, err := s.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	// Aqui seria implementada a lógica para adicionar o token à lista negra
	// Por exemplo, armazenar em um banco de dados ou cache

	return nil
}

// IsTokenRevoked verifica se um token está revogado
func (s *DefaultJWTService) IsTokenRevoked(tokenString string) (bool, error) {
	// Aqui seria implementada a lógica para verificar se o token está na lista negra
	// Por exemplo, consultar um banco de dados ou cache

	// Por enquanto, retornamos false (não revogado)
	return false, nil
}