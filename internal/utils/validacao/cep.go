package validacao

import (
	"regexp"
)

// ValidarCEP valida um CEP brasileiro
// Recebe o CEP a ser validado
// Retorna true se o CEP for válido, false caso contrário
func ValidarCEP(cep string) bool {
	// Remover caracteres não numéricos
	re := regexp.MustCompile(`[^0-9]`)
	cep = re.ReplaceAllString(cep, "")

	// Verificar tamanho (8 dígitos)
	return len(cep) == 8
}

// FormatarCEP formata um CEP no padrão XXXXX-XXX
// Recebe o CEP a ser formatado
// Retorna o CEP formatado
func FormatarCEP(cep string) string {
	// Remover caracteres não numéricos
	re := regexp.MustCompile(`[^0-9]`)
	cep = re.ReplaceAllString(cep, "")

	// Verificar tamanho
	if len(cep) != 8 {
		return cep
	}

	// Formatar CEP
	return cep[0:5] + "-" + cep[5:8]
}