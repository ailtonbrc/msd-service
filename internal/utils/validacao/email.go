package validacao

import (
	"regexp"
	"strings"
)

// Expressão regular para validar emails
// Esta regex verifica:
// - Parte local (antes do @) contém caracteres válidos
// - Domínio (após o @) tem formato válido
// - TLD (parte final após o último ponto) tem pelo menos 2 caracteres
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidarEmail verifica se um email está em formato válido
// Retorna true se o email for válido, false caso contrário
func ValidarEmail(email string) bool {
	// Verificar se o email está vazio
	if email == "" {
		return false
	}

	// Remover espaços em branco no início e fim
	email = strings.TrimSpace(email)

	// Verificar tamanho mínimo e máximo
	if len(email) < 6 || len(email) > 254 {
		return false
	}

	// Verificar se contém apenas um @
	if strings.Count(email, "@") != 1 {
		return false
	}

	// Verificar formato usando regex
	if !emailRegex.MatchString(email) {
		return false
	}

	// Verificações adicionais
	parts := strings.Split(email, "@")
	localPart := parts[0]
	domainPart := parts[1]

	// Verificar se a parte local não começa ou termina com ponto
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return false
	}

	// Verificar se não há pontos consecutivos
	if strings.Contains(localPart, "..") || strings.Contains(domainPart, "..") {
		return false
	}

	// Verificar se o domínio não começa ou termina com hífen
	if strings.HasPrefix(domainPart, "-") || strings.HasSuffix(domainPart, "-") {
		return false
	}

	return true
}

// FormatarEmail padroniza um email, removendo espaços e convertendo para minúsculas
// Retorna o email formatado ou uma string vazia se o email for inválido
func FormatarEmail(email string) string {
	// Remover espaços e converter para minúsculas
	email = strings.ToLower(strings.TrimSpace(email))
	
	// Verificar se é válido
	if !ValidarEmail(email) {
		return ""
	}
	
	return email
}