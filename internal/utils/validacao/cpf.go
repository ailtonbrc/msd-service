package validacao

import (
	"regexp"
	"strconv"
)

// ValidarCPF valida um CPF brasileiro
// Recebe o CPF a ser validado
// Retorna true se o CPF for válido, false caso contrário
func ValidarCPF(cpf string) bool {
	// Remover caracteres não numéricos
	re := regexp.MustCompile(`[^0-9]`)
	cpf = re.ReplaceAllString(cpf, "")

	// Verificar tamanho
	if len(cpf) != 11 {
		return false
	}

	// Verificar se todos os dígitos são iguais
	if allDigitsEqual(cpf) {
		return false
	}

	// Calcular primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	if remainder < 2 {
		remainder = 0
	} else {
		remainder = 11 - remainder
	}

	// Verificar primeiro dígito
	firstDigit, _ := strconv.Atoi(string(cpf[9]))
	if remainder != firstDigit {
		return false
	}

	// Calcular segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	if remainder < 2 {
		remainder = 0
	} else {
		remainder = 11 - remainder
	}

	// Verificar segundo dígito
	secondDigit, _ := strconv.Atoi(string(cpf[10]))
	return remainder == secondDigit
}

// allDigitsEqual verifica se todos os dígitos de uma string são iguais
func allDigitsEqual(s string) bool {
	if len(s) == 0 {
		return true
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}

// FormatarCPF formata um CPF no padrão XXX.XXX.XXX-XX
// Recebe o CPF a ser formatado
// Retorna o CPF formatado
func FormatarCPF(cpf string) string {
	// Remover caracteres não numéricos
	re := regexp.MustCompile(`[^0-9]`)
	cpf = re.ReplaceAllString(cpf, "")

	// Verificar tamanho
	if len(cpf) != 11 {
		return cpf
	}

	// Formatar CPF
	return cpf[0:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:11]
}