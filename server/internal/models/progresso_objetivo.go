package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProgressoObjetivo representa o progresso de um objetivo terapêutico
type ProgressoObjetivo struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ObjetivoID  uuid.UUID      `gorm:"type:uuid;not null" json:"objetivo_id"`
	Data        time.Time      `gorm:"not null" json:"data"`
	Nota        int            `gorm:"not null" json:"nota"`
	Observacoes string         `gorm:"type:text" json:"observacoes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (ProgressoObjetivo) TableName() string {
	return "progresso_objetivo"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (p *ProgressoObjetivo) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
