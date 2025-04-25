package validacao

import (
	"regexp"
	"strconv"
)

// ValidarCPF verifica se o CPF é válido
func ValidarCPF(cpf string) bool {
	// Se o CPF estiver vazio, consideramos inválido
	if cpf == "" {
		return false
	}

	// Remove caracteres não numéricos
	cpf = LimparCPF(cpf)

	// Verifica se o CPF tem 11 dígitos
	if len(cpf) != 11 {
		return false
	}

	// Verifica se todos os dígitos são iguais
	if todosDigitosIguais(cpf) {
		return false
	}

	// Calcula e verifica o primeiro dígito verificador
	d1 := calcularDigitoVerificador(cpf, 9)
	if d1 != int(cpf[9]-'0') {
		return false
	}

	// Calcula e verifica o segundo dígito verificador
	d2 := calcularDigitoVerificador(cpf, 10)
	if d2 != int(cpf[10]-'0') {
		return false
	}

	return true
}

// LimparCPF remove caracteres não numéricos do CPF
func LimparCPF(cpf string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(cpf, "")
}

// todosDigitosIguais verifica se todos os dígitos do CPF são iguais
func todosDigitosIguais(cpf string) bool {
	first := cpf[0]
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != first {
			return false
		}
	}
	return true
}

// calcularDigitoVerificador calcula o dígito verificador do CPF
func calcularDigitoVerificador(cpf string, pos int) int {
	var soma int
	var peso int = pos + 1

	// Soma os produtos dos dígitos pelos pesos
	for i := 0; i < pos; i++ {
		digito, _ := strconv.Atoi(string(cpf[i]))
		soma += digito * peso
		peso--
	}

	// Calcula o resto da divisão por 11
	resto := soma % 11
	
	// Se o resto for menor que 2, o dígito é 0, senão é 11 - resto
	if resto < 2 {
		return 0
	}
	return 11 - resto
}

// FormatarCPF formata o CPF no padrão XXX.XXX.XXX-XX
func FormatarCPF(cpf string) string {
	cpf = LimparCPF(cpf)
	if len(cpf) != 11 {
		return cpf // Retorna o original se não tiver 11 dígitos
	}
	
	return cpf[0:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:11]
}