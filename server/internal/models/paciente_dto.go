package models

import (
	"time"

	"github.com/google/uuid"
)

// CreatePacienteRequest representa os dados necessários para criar um novo paciente
type CreatePacienteRequest struct {
	Nome           string    `json:"nome" binding:"required" example:"João Silva"`
	DataNascimento time.Time `json:"data_nascimento" binding:"required" example:"2015-01-01T00:00:00Z"`
	GrauTEA        GrauTEA   `json:"grau_tea" binding:"required,oneof=Leve Moderado Severo" example:"Leve"`
	ResponsavelID  uuid.UUID `json:"responsavel_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Observacoes    string    `json:"observacoes" example:"Paciente apresenta dificuldades de comunicação"`
}

// UpdatePacienteRequest representa os dados que podem ser atualizados em um paciente
type UpdatePacienteRequest struct {
	Nome           *string    `json:"nome" example:"João Silva"`
	DataNascimento *time.Time `json:"data_nascimento" example:"2015-01-01T00:00:00Z"`
	GrauTEA        *GrauTEA   `json:"grau_tea" binding:"omitempty,oneof=Leve Moderado Severo" example:"Moderado"`
	ResponsavelID  *uuid.UUID `json:"responsavel_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Observacoes    *string    `json:"observacoes" example:"Paciente apresenta melhoras na comunicação"`
}

// PacienteResponse representa os dados de um paciente que são retornados pela API
type PacienteResponse struct {
	ID             uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Nome           string    `json:"nome" example:"João Silva"`
	DataNascimento time.Time `json:"data_nascimento" example:"2015-01-01T00:00:00Z"`
	GrauTEA        GrauTEA   `json:"grau_tea" example:"Leve"`
	ResponsavelID  uuid.UUID `json:"responsavel_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Observacoes    string    `json:"observacoes" example:"Paciente apresenta dificuldades de comunicação"`
	CreatedAt      time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ToPaciente converte um CreatePacienteRequest para um modelo Paciente
func (r *CreatePacienteRequest) ToPaciente() *Paciente {
	return &Paciente{
		Nome:           r.Nome,
		DataNascimento: r.DataNascimento,
		GrauTEA:        r.GrauTEA,
		ResponsavelID:  r.ResponsavelID,
		Observacoes:    r.Observacoes,
	}
}

// ToResponse converte um modelo Paciente para um PacienteResponse
func (p *Paciente) ToResponse() *PacienteResponse {
	return &PacienteResponse{
		ID:             p.ID,
		Nome:           p.Nome,
		DataNascimento: p.DataNascimento,
		GrauTEA:        p.GrauTEA,
		ResponsavelID:  p.ResponsavelID,
		Observacoes:    p.Observacoes,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

// ApplyUpdates aplica as atualizações de um UpdatePacienteRequest a um modelo Paciente
func (p *Paciente) ApplyUpdates(updates *UpdatePacienteRequest) {
	if updates.Nome != nil {
		p.Nome = *updates.Nome
	}
	if updates.DataNascimento != nil {
		p.DataNascimento = *updates.DataNascimento
	}
	if updates.GrauTEA != nil {
		p.GrauTEA = *updates.GrauTEA
	}
	if updates.ResponsavelID != nil {
		p.ResponsavelID = *updates.ResponsavelID
	}
	if updates.Observacoes != nil {
		p.Observacoes = *updates.Observacoes
	}
}
