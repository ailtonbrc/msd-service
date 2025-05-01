package validacao

import (
	"regexp"
)

// ValidarTelefone valida um número de telefone brasileiro
// Recebe o telefone a ser validado
// Retorna true se o telefone for válido, false caso contrário
func ValidarTelefone(telefone string) bool {
	// Remover caracteres não numéricos
	re := regexp.MustCompile(`[^0-9]`)
	telefone = re.ReplaceAllString(telefone, "")

	// Verificar tamanho (8 a 11 dígitos)
	// 8 dígitos: telefone fixo sem DDD
	// 9 dígitos: celular sem DDD
	// 10 dígitos: telefone fixo com DDD
	// 11 dígitos: celular com DDD
	if len(telefone) < 8 || len(telefone) > 11 {
		return false
	}

	// Verificar se é um número de celular (9 dígitos)
	if len(telefone) == 9 {
		// O primeiro dígito deve ser 9
		return telefone[0] == '9'
	}

	// Verificar se é um número de celular com DDD (11 dígitos)
	if len(telefone) == 11 {
		// O terceiro dígito deve ser 9
		return telefone[2] == '9'
	}

	return true
}

// FormatarTelefone formata um número de telefone brasileiro
// Recebe o telefone a ser formatado
// Retorna o telefone formatado
func FormatarTelefone(telefone string) string {
	// Remover caracteres não numéricos
	re := regexp.MustCompile(`[^0-9]`)
	telefone = re.ReplaceAllString(telefone, "")

	// Formatar de acordo com o tamanho
	switch len(telefone) {
	case 8: // Telefone fixo sem DDD
		return telefone[0:4] + "-" + telefone[4:8]
	case 9: // Celular sem DDD
		return telefone[0:5] + "-" + telefone[5:9]
	case 10: // Telefone fixo com DDD
		return "(" + telefone[0:2] + ") " + telefone[2:6] + "-" + telefone[6:10]
	case 11: // Celular com DDD
		return "(" + telefone[0:2] + ") " + telefone[2:7] + "-" + telefone[7:11]
	default:
		return telefone
	}
}