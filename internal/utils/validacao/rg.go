package validacao

import (
	"regexp"
)

// ValidarRG valida um RG brasileiro
// Recebe o RG a ser validado
// Retorna true se o RG for válido, false caso contrário
func ValidarRG(rg string) bool {
	// Remover caracteres não alfanuméricos
	re := regexp.MustCompile(`[^0-9A-Za-z]`)
	rg = re.ReplaceAllString(rg, "")

	// Verificar tamanho (geralmente entre 5 e 14 caracteres)
	return len(rg) >= 5 && len(rg) <= 14
}

// FormatarRG formata um RG no padrão XX.XXX.XXX-X
// Recebe o RG a ser formatado
// Retorna o RG formatado
func FormatarRG(rg string) string {
	// Remover caracteres não alfanuméricos
	re := regexp.MustCompile(`[^0-9A-Za-z]`)
	rg = re.ReplaceAllString(rg, "")

	// Se o RG for muito curto, retornar sem formatação
	if len(rg) < 8 {
		return rg
	}

	// Formatar RG (considerando o formato mais comum)
	return rg[0:2] + "." + rg[2:5] + "." + rg[5:8] + "-" + rg[8:]
}