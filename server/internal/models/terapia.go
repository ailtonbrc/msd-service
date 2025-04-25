package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Terapia representa um tipo de terapia oferecida pela clínica
type Terapia struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nome      string         `gorm:"size:100;not null" json:"nome"`
	Descricao string         `gorm:"type:text" json:"descricao"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (Terapia) TableName() string {
	return "terapias"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (t *Terapia) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
