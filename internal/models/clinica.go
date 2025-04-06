package models

import "time"

type Clinica struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Nome      string    `json:"nome"`
    CNPJ      string    `gorm:"unique" json:"cnpj"`
    Endereco  string    `json:"endereco"`
    Telefone  string    `json:"telefone"`
    CriadoEm  time.Time `gorm:"autoCreateTime" json:"criado_em"`
}

func (Clinica) TableName() string {
    return "clinicas"
}