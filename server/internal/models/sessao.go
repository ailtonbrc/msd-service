package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StatusSessao representa o status de uma sessão
type StatusSessao string

const (
	StatusSessaoPlanejada StatusSessao = "planejada"
	StatusSessaoRealizada StatusSessao = "realizada"
	StatusSessaoCancelada StatusSessao = "cancelada"
)

// Sessao representa uma sessão de terapia
type Sessao struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PacienteID     uuid.UUID      `gorm:"type:uuid;not null" json:"paciente_id"`
	TerapeutaID    uuid.UUID      `gorm:"type:uuid;not null" json:"terapeuta_id"`
	TerapiaID      uuid.UUID      `gorm:"type:uuid;not null" json:"terapia_id"`
	Data           time.Time      `gorm:"not null" json:"data"`
	DuracaoMinutos int            `gorm:"not null" json:"duracao_minutos"`
	Status         StatusSessao   `gorm:"type:varchar(20);not null" json:"status"`
	ResumoSessao   string         `gorm:"type:text" json:"resumo_sessao"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (Sessao) TableName() string {
	return "sessoes"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (s *Sessao) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
