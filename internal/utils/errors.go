package utils

import (
	"errors"
)

// Erros comuns
var (
	ErrNotFound           = errors.New("registro não encontrado")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUnauthorized       = errors.New("não autorizado")
	ErrForbidden          = errors.New("acesso negado")
	ErrInvalidInput       = errors.New("dados de entrada inválidos")
	ErrDuplicateEntry     = errors.New("registro duplicado")
	ErrInternalServer     = errors.New("erro interno do servidor")
)