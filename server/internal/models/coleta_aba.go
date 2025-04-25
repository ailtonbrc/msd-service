package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResultadoColeta representa o resultado de uma coleta ABA
type ResultadoColeta string

const (
	ResultadoColetaAcerto ResultadoColeta = "acerto"
	ResultadoColetaErro   ResultadoColeta = "erro"
	ResultadoColetaAjuda  ResultadoColeta = "ajuda"
)

// ColetaABA representa uma coleta de dados em uma sessão ABA
type ColetaABA struct {
	ID               uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	EtapaProgramaID  uuid.UUID       `gorm:"type:uuid;not null" json:"etapa_programa_id"`
	SessaoID         uuid.UUID       `gorm:"type:uuid;not null" json:"sessao_id"`
	Resultado        ResultadoColeta `gorm:"type:varchar(20);not null" json:"resultado"`
	PromptUtilizadoID uuid.UUID      `gorm:"type:uuid" json:"prompt_utilizado_id"`
	ReforcoUtilizado string          `gorm:"size:100" json:"reforco_utilizado"`
	Observacoes      string          `gorm:"type:text" json:"observacoes"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName especifica o nome da tabela no banco de dados
func (ColetaABA) TableName() string {
	return "coletas_aba"
}

// BeforeCreate é um hook do GORM que é executado antes de criar um registro
func (c *ColetaABA) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
