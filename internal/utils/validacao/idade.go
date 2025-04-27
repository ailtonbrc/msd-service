package validacao

import (
	"time"
)

// CalculateAge calcula a idade com base na data de nascimento
func Calculaidade(dataAniversario time.Time) int {
	if dataAniversario.IsZero() {
		return 0
	}
	
	hoje := time.Now()
	idade := hoje.Year() - dataAniversario.Year()
	
	// Ajusta a idade se ainda não fez aniversário este ano
	if hoje.Month() < dataAniversario.Month() || 
	   (hoje.Month() == dataAniversario.Month() && hoje.Day() < dataAniversario.Day()) {
		idade--
	}
	
	return idade
}