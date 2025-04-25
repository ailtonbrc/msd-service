package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EtapaPrograma representa uma etapa de um programa ABA
type EtapaPrograma struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProgramaID      uuid.UUID      `gorm:"type:uuid;not null" json:"programa_id"`
	Descricao       string         `gorm:"type:text;not null" json:"descricao"`
	Ordem           int            `gorm:"not null" json:"ordem"`
	CriterioSucesso string         `gorm:"type:text" json:"criterio_sucesso"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (EtapaPrograma) TableName() string {
	return "etapas_programa"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (e *EtapaPrograma) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}
