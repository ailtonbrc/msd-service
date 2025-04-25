package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TipoPrompt representa um tipo de prompt utilizado nas terapias ABA
type TipoPrompt struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Tipo      string         `gorm:"size:50;not null" json:"tipo"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (TipoPrompt) TableName() string {
	return "tipos_prompt"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (t *TipoPrompt) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
