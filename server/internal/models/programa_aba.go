package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StatusPrograma representa o status de um programa ABA
type StatusPrograma string

const (
	StatusProgramaAtivo     StatusPrograma = "ativo"
	StatusProgramaPausado   StatusPrograma = "pausado"
	StatusProgramaFinalizado StatusPrograma = "finalizado"
)

// ProgramaABA representa um programa ABA para um paciente
type ProgramaABA struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nome       string         `gorm:"size:100;not null" json:"nome"`
	Descricao  string         `gorm:"type:text" json:"descricao"`
	PacienteID uuid.UUID      `gorm:"type:uuid;not null" json:"paciente_id"`
	DataInicio time.Time      `gorm:"not null" json:"data_inicio"`
	DataFim    time.Time      `json:"data_fim"`
	Status     StatusPrograma `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (ProgramaABA) TableName() string {
	return "programas_aba"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (p *ProgramaABA) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
