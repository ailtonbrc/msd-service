// internal/auth/permissoes.go
package auth

import (
	"context"
	"strings"
)

// Chave para armazenar os claims do usuário no contexto
type contextKey string
const UserClaimsKey contextKey = "user_claims"

// GetUserFromContext extrai os claims do usuário do contexto
func GetUserFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*UserClaims)
	return claims, ok
}

// HasPermission verifica se o usuário tem a permissão especificada
func HasPermission(claims *UserClaims, permission string) bool {
	if claims == nil {
		return false
	}

	// Administradores têm todas as permissões
	if HasRole(claims, "admin") {
		return true
	}

	// Verificar permissões específicas
	for _, p := range claims.Permissions {
		// Verificar permissão exata
		if p == permission {
			return true
		}

		// Verificar permissão curinga (ex: "pacientes:*")
		if strings.HasSuffix(p, ":*") {
			prefix := strings.TrimSuffix(p, ":*")
			if strings.HasPrefix(permission, prefix+":") {
				return true
			}
		}
	}

	return false
}

// HasRole verifica se o usuário tem o papel especificado
func HasRole(claims *UserClaims, role string) bool {
	if claims == nil {
		return false
	}

	for _, r := range claims.Roles {
		if r == role {
			return true
		}
	}

	return false
}

// HasScope verifica se o usuário tem o escopo especificado
func HasScope(claims *UserClaims, scope string) bool {
	if claims == nil {
		return false
	}

	for _, s := range claims.Scopes {
		if s == scope {
			return true
		}

		// Verificar escopo curinga (ex: "api:*")
		if strings.HasSuffix(s, ":*") {
			prefix := strings.TrimSuffix(s, ":*")
			if strings.HasPrefix(scope, prefix+":") {
				return true
			}
		}
	}

	return false
}