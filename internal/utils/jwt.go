package utils

import (
	"errors"
	"time"

	"clinica_server/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims representa os claims do token JWT
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Perfil   string `json:"perfil"`
	jwt.RegisteredClaims
}

// GenerateAccessToken gera um novo token JWT de acesso
func GenerateAccessToken(userID uint, username string, perfil string, cfg *config.Config) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Perfil:   perfil,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.AccessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "clinica_server",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// GenerateRefreshToken gera um novo token JWT de refresh
func GenerateRefreshToken(userID uint, username string, cfg *config.Config) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshTokenExp)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "clinica_server",
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ValidateToken valida um token JWT
func ValidateToken(tokenString string, cfg *config.Config) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
