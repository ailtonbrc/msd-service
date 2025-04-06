package models

import "time"

type Paciente struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Nome                string    `json:"nome"`
	DataNascimento      time.Time `json:"data_nascimento"`
	Sexo                string    `json:"sexo"`
	Responsavel         string    `json:"responsavel"`
	TelefoneResponsavel string    `json:"telefone_responsavel"`
	EmailResponsavel    string    `json:"email_responsavel"`
	Endereco            string    `json:"endereco"`
	Diagnostico         string    `json:"diagnostico"` // Ex: TEA leve, moderado, severo
	DataDiagnostico     *time.Time `json:"data_diagnostico"`
	LaudoMedico         string    `json:"laudo_medico"` // Referência ou resumo do laudo
	Observacoes         string    `json:"observacoes"`
	ClinicaID           *uint     `json:"clinica_id"`
	CriadoEm            time.Time `gorm:"autoCreateTime" json:"criado_em"`
	AtualizadoEm        time.Time `gorm:"autoUpdateTime" json:"atualizado_em"`
}

func (Paciente) TableName() string {
	return "pacientes"
}