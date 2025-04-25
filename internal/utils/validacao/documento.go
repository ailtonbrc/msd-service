package validacao

// Documento define uma interface para documentos que podem ser validados
type Documento interface {
	Validar() bool
	Formatar() string
	Limpar() string
}

// CPF implementa a interface Documento
type CPF string

func (c CPF) Validar() bool {
	return ValidarCPF(string(c))
}

func (c CPF) Formatar() string {
	return FormatarCPF(string(c))
}

func (c CPF) Limpar() string {
	return LimparCPF(string(c))
}

// Criar um novo CPF a partir de uma string
func NovoCPF(cpf string) CPF {
	return CPF(cpf)
}