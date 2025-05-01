package validacao

import (
	"regexp"
	"strings"
)

// ValidarEmail valida um endereço de email
// Recebe o email a ser validado
// Retorna true se o email for válido, false caso contrário
func ValidarEmail(email string) bool {
	// Verificar se o email está vazio
	if email == "" {
		return false
	}

	// Remover espaços em branco
	email = strings.TrimSpace(email)

	// Verificar formato do email usando expressão regular
	// Esta regex verifica o formato básico de um email: usuario@dominio.com
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// NormalizarEmail normaliza um endereço de email (converte para minúsculas)
// Recebe o email a ser normalizado
// Retorna o email normalizado
func NormalizarEmail(email string) string {
	// Remover espaços em branco e converter para minúsculas
	return strings.ToLower(strings.TrimSpace(email))
}