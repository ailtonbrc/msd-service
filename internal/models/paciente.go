package models

import (
	"time"

	"gorm.io/gorm"
)

// Paciente representa um paciente no sistema
type Paciente struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	Nome                string         `json:"nome" gorm:"size:100;not null"`
	DataNascimento      time.Time      `json:"data_nascimento"`
	Genero              string         `json:"genero" gorm:"size:20"`
	CPF                 string         `json:"cpf" gorm:"size:14;uniqueIndex"`
	RG                  string         `json:"rg" gorm:"size:20"`
	Diagnostico         string         `json:"diagnostico" gorm:"size:255"`
	Telefone            string         `json:"telefone" gorm:"size:20"`
	Email               string         `json:"email" gorm:"size:100"`
	Endereco            string         `json:"endereco" gorm:"size:255"`
	Cidade              string         `json:"cidade" gorm:"size:100"`
	Estado              string         `json:"estado" gorm:"size:2"`
	CEP                 string         `json:"cep" gorm:"size:10"`
	NomeResponsavel     string         `json:"nome_responsavel" gorm:"size:100"`
	TelefoneResponsavel string         `json:"telefone_responsavel" gorm:"size:20"`
	EmailResponsavel    string         `json:"email_responsavel" gorm:"size:100"`
	Observacoes         string         `json:"observacoes" gorm:"type:text"`
	Alergias            string         `json:"alergias" gorm:"type:text"`
	Medicacoes          string         `json:"medicacoes" gorm:"type:text"`
	CriadoEm            time.Time      `json:"criado_em" gorm:"autoCreateTime"`
	AtualizadoEm        time.Time      `json:"atualizado_em" gorm:"autoUpdateTime"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
	CriadoPor           uint           `json:"criado_por" gorm:"not null"` // ID do usuário que criou
	AtualizadoPor       uint           `json:"atualizado_por"`             // ID do usuário que atualizou
}

// TableName define o nome da tabela no banco de dados
func (Paciente) TableName() string {
	return "pacientes"
}

// PacienteFiltro define os filtros para busca avançada de pacientes
type PacienteFiltro struct {
	Nome        string
	CPF         string
	Telefone    string
	Diagnostico string
	Page        int
	PageSize    int
}

// PaginatedResponse representa uma resposta paginada
type PaginatedResponse struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Data     interface{} `json:"data"`
}