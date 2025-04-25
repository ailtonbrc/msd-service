package models

import (
	"gorm.io/gorm"
)

type Clinica struct {
	gorm.Model

	Nome     string `gorm:"not null" json:"nome"`
	CNPJ     string `gorm:"unique;not null" json:"cnpj"`
	Endereco string `json:"endereco"`
	Telefone string `json:"telefone"`
}

// TableName especifica o nome da tabela
func (Clinica) TableName() string {
	return "clinica"
}
