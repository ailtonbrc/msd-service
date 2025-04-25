package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StatusObjetivo representa o status de um objetivo terapêutico
type StatusObjetivo string

const (
	StatusObjetivoEmProgresso StatusObjetivo = "em progresso"
	StatusObjetivoConcluido   StatusObjetivo = "concluido"
	StatusObjetivoSuspenso    StatusObjetivo = "suspenso"
)

// ObjetivoTerapeutico representa um objetivo terapêutico para um paciente
type ObjetivoTerapeutico struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PacienteID uuid.UUID      `gorm:"type:uuid;not null" json:"paciente_id"`
	Descricao  string         `gorm:"type:text;not null" json:"descricao"`
	DataInicio time.Time      `gorm:"not null" json:"data_inicio"`
	DataFim    time.Time      `json:"data_fim"`
	Status     StatusObjetivo `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (ObjetivoTerapeutico) TableName() string {
	return "objetivos_terapeuticos"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (o *ObjetivoTerapeutico) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
