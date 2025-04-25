package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegistroComportamento representa um registro de ocorrência de comportamento
type RegistroComportamento struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ComportamentoID uuid.UUID      `gorm:"type:uuid;not null" json:"comportamento_id"`
	DataHora        time.Time      `gorm:"not null" json:"data_hora"`
	Valor           float64        `gorm:"not null" json:"valor"`
	Contexto        string         `gorm:"type:text" json:"contexto"`
	Consequencia    string         `gorm:"type:text" json:"consequencia"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (RegistroComportamento) TableName() string {
	return "registros_comportamento"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (r *RegistroComportamento) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
