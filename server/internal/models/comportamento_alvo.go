package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TipoComportamento representa o tipo de um comportamento alvo
type TipoComportamento string

const (
	TipoComportamentoAdequado   TipoComportamento = "adequado"
	TipoComportamentoInadequado TipoComportamento = "inadequado"
)

// MetodoRegistro representa o método de registro de um comportamento
type MetodoRegistro string

const (
	MetodoRegistroFrequencia MetodoRegistro = "frequencia"
	MetodoRegistroDuracao    MetodoRegistro = "duracao"
	MetodoRegistroIntensidade MetodoRegistro = "intensidade"
	MetodoRegistroIntervalo   MetodoRegistro = "intervalo"
)

// ComportamentoAlvo representa um comportamento alvo para monitoramento
type ComportamentoAlvo struct {
	ID             uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PacienteID     uuid.UUID        `gorm:"type:uuid;not null" json:"paciente_id"`
	Descricao      string           `gorm:"type:text;not null" json:"descricao"`
	Tipo           TipoComportamento `gorm:"type:varchar(20);not null" json:"tipo"`
	MetodoRegistro MetodoRegistro    `gorm:"type:varchar(20);not null" json:"metodo_registro"`
	DataInicio     time.Time         `gorm:"not null" json:"data_inicio"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      gorm.DeletedAt    `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (ComportamentoAlvo) TableName() string {
	return "comportamentos_alvo"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (c *ComportamentoAlvo) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
