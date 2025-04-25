package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GrauTEA representa o grau do Transtorno do Espectro Autista
type GrauTEA string

const (
	GrauTEALeve     GrauTEA = "Leve"
	GrauTEAModerado GrauTEA = "Moderado"
	GrauTEASevero   GrauTEA = "Severo"
)

// Paciente representa um paciente da clínica
type Paciente struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nome           string         `gorm:"size:100;not null" json:"nome"`
	DataNascimento time.Time      `gorm:"not null" json:"data_nascimento"`
	GrauTEA        GrauTEA        `gorm:"type:varchar(20);not null" json:"grau_tea"`
	ResponsavelID  uuid.UUID      `gorm:"type:uuid" json:"responsavel_id"`
	Observacoes    string         `gorm:"type:text" json:"observacoes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (Paciente) TableName() string {
	return "pacientes"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (p *Paciente) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
