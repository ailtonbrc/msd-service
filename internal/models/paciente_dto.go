package models

import (
	"time"
)

// CreatePacienteRequest representa os dados para criação de um paciente
type CreatePacienteRequest struct {
	Nome                string    `json:"nome" binding:"required,max=100"`
	DataNascimento      time.Time `json:"data_nascimento" binding:"required"`
	Genero              string    `json:"genero" binding:"max=20"`
	CPF                 string    `json:"cpf" binding:"omitempty,len=14"`
	RG                  string    `json:"rg" binding:"max=20"`
	Diagnostico         string    `json:"diagnostico" binding:"max=255"`
	Telefone            string    `json:"telefone" binding:"required,max=20"`
	Email               string    `json:"email" binding:"omitempty,email,max=100"`
	Endereco            string    `json:"endereco" binding:"max=255"`
	Cidade              string    `json:"cidade" binding:"max=100"`
	Estado              string    `json:"estado" binding:"max=2"`
	CEP                 string    `json:"cep" binding:"max=10"`
	NomeResponsavel     string    `json:"nome_responsavel" binding:"omitempty,max=100"`
	TelefoneResponsavel string    `json:"telefone_responsavel" binding:"omitempty,max=20"`
	EmailResponsavel    string    `json:"email_responsavel" binding:"omitempty,email,max=100"`
	Observacoes         string    `json:"observacoes"`
	Alergias            string    `json:"alergias"`
	Medicacoes          string    `json:"medicacoes"`
}

// UpdatePacienteRequest representa os dados para atualização de um paciente
type UpdatePacienteRequest struct {
	Nome                string    `json:"nome" binding:"omitempty,max=100"`
	DataNascimento      time.Time `json:"data_nascimento"`
	Genero              string    `json:"genero" binding:"max=20"`
	CPF                 string    `json:"cpf" binding:"omitempty,len=14"`
	RG                  string    `json:"rg" binding:"max=20"`
	Diagnostico         string    `json:"diagnostico" binding:"max=255"`
	Telefone            string    `json:"telefone" binding:"omitempty,max=20"`
	Email               string    `json:"email" binding:"omitempty,email,max=100"`
	Endereco            string    `json:"endereco" binding:"max=255"`
	Cidade              string    `json:"cidade" binding:"max=100"`
	Estado              string    `json:"estado" binding:"max=2"`
	CEP                 string    `json:"cep" binding:"max=10"`
	NomeResponsavel     string    `json:"nome_responsavel" binding:"omitempty,max=100"`
	TelefoneResponsavel string    `json:"telefone_responsavel" binding:"omitempty,max=20"`
	EmailResponsavel    string    `json:"email_responsavel" binding:"omitempty,email,max=100"`
	Observacoes         string    `json:"observacoes"`
	Alergias            string    `json:"alergias"`
	Medicacoes          string    `json:"medicacoes"`
}

// PacienteResponse representa a resposta com dados de um paciente
type PacienteResponse struct {
	ID                  uint      `json:"id"`
	Nome                string    `json:"nome"`
	DataNascimento      time.Time `json:"data_nascimento"`
	Genero              string    `json:"genero"`
	CPF                 string    `json:"cpf"`
	RG                  string    `json:"rg,omitempty"`
	Diagnostico         string    `json:"diagnostico,omitempty"`
	Telefone            string    `json:"telefone"`
	Email               string    `json:"email,omitempty"`
	Endereco            string    `json:"endereco,omitempty"`
	Cidade              string    `json:"cidade,omitempty"`
	Estado              string    `json:"estado,omitempty"`
	CEP                 string    `json:"cep,omitempty"`
	NomeResponsavel     string    `json:"nome_responsavel,omitempty"`
	TelefoneResponsavel string    `json:"telefone_responsavel,omitempty"`
	EmailResponsavel    string    `json:"email_responsavel,omitempty"`
	Observacoes         string    `json:"observacoes,omitempty"`
	Alergias            string    `json:"alergias,omitempty"`
	Medicacoes          string    `json:"medicacoes,omitempty"`
	CriadoEm            time.Time `json:"criado_em"`
	AtualizadoEm        time.Time `json:"atualizado_em"`
	CriadoPor           uint      `json:"criado_por"`
	AtualizadoPor       uint      `json:"atualizado_por,omitempty"`
}

// PacienteDTOConverter define a interface para conversão entre DTOs e modelos de pacientes
type PacienteDTOConverter interface {
	CreateRequestToModel(req *CreatePacienteRequest, userID uint) *Paciente
	UpdateRequestToModel(req *UpdatePacienteRequest, existingPaciente *Paciente, userID uint) *Paciente
	ModelToResponse(paciente *Paciente) *PacienteResponse
}

// DefaultPacienteDTOConverter implementa PacienteDTOConverter
type DefaultPacienteDTOConverter struct{}

// NewPacienteDTOConverter cria uma nova instância de PacienteDTOConverter
func NewPacienteDTOConverter() PacienteDTOConverter {
	return &DefaultPacienteDTOConverter{}
}

// CreateRequestToModel converte um CreatePacienteRequest para um modelo Paciente
func (c *DefaultPacienteDTOConverter) CreateRequestToModel(req *CreatePacienteRequest, userID uint) *Paciente {
	return &Paciente{
		Nome:                req.Nome,
		DataNascimento:      req.DataNascimento,
		Genero:              req.Genero,
		CPF:                 req.CPF,
		RG:                  req.RG,
		Diagnostico:         req.Diagnostico,
		Telefone:            req.Telefone,
		Email:               req.Email,
		Endereco:            req.Endereco,
		Cidade:              req.Cidade,
		Estado:              req.Estado,
		CEP:                 req.CEP,
		NomeResponsavel:     req.NomeResponsavel,
		TelefoneResponsavel: req.TelefoneResponsavel,
		EmailResponsavel:    req.EmailResponsavel,
		Observacoes:         req.Observacoes,
		Alergias:            req.Alergias,
		Medicacoes:          req.Medicacoes,
		CriadoEm:            time.Now(),
		AtualizadoEm:        time.Now(),
		CriadoPor:           userID,
		AtualizadoPor:       userID,
	}
}

// UpdateRequestToModel converte um UpdatePacienteRequest para um modelo Paciente
func (c *DefaultPacienteDTOConverter) UpdateRequestToModel(req *UpdatePacienteRequest, existingPaciente *Paciente, userID uint) *Paciente {
	// Criar uma cópia do paciente existente
	paciente := *existingPaciente
	
	// Atualizar apenas os campos fornecidos
	if req.Nome != "" {
		paciente.Nome = req.Nome
	}
	
	if !req.DataNascimento.IsZero() {
		paciente.DataNascimento = req.DataNascimento
	}
	
	if req.Genero != "" {
		paciente.Genero = req.Genero
	}
	
	if req.CPF != "" {
		paciente.CPF = req.CPF
	}
	
	if req.RG != "" {
		paciente.RG = req.RG
	}
	
	if req.Diagnostico != "" {
		paciente.Diagnostico = req.Diagnostico
	}
	
	if req.Telefone != "" {
		paciente.Telefone = req.Telefone
	}
	
	if req.Email != "" {
		paciente.Email = req.Email
	}
	
	if req.Endereco != "" {
		paciente.Endereco = req.Endereco
	}
	
	if req.Cidade != "" {
		paciente.Cidade = req.Cidade
	}
	
	if req.Estado != "" {
		paciente.Estado = req.Estado
	}
	
	if req.CEP != "" {
		paciente.CEP = req.CEP
	}
	
	if req.NomeResponsavel != "" {
		paciente.NomeResponsavel = req.NomeResponsavel
	}
	
	if req.TelefoneResponsavel != "" {
		paciente.TelefoneResponsavel = req.TelefoneResponsavel
	}
	
	if req.EmailResponsavel != "" {
		paciente.EmailResponsavel = req.EmailResponsavel
	}
	
	// Estes campos sempre são atualizados
	paciente.Observacoes = req.Observacoes
	paciente.Alergias = req.Alergias
	paciente.Medicacoes = req.Medicacoes
	paciente.AtualizadoEm = time.Now()
	paciente.AtualizadoPor = userID
	
	return &paciente
}

// ModelToResponse converte um modelo Paciente para um PacienteResponse
func (c *DefaultPacienteDTOConverter) ModelToResponse(paciente *Paciente) *PacienteResponse {
	return &PacienteResponse{
		ID:                  paciente.ID,
		Nome:                paciente.Nome,
		DataNascimento:      paciente.DataNascimento,
		Genero:              paciente.Genero,
		CPF:                 paciente.CPF,
		RG:                  paciente.RG,
		Diagnostico:         paciente.Diagnostico,
		Telefone:            paciente.Telefone,
		Email:               paciente.Email,
		Endereco:            paciente.Endereco,
		Cidade:              paciente.Cidade,
		Estado:              paciente.Estado,
		CEP:                 paciente.CEP,
		NomeResponsavel:     paciente.NomeResponsavel,
		TelefoneResponsavel: paciente.TelefoneResponsavel,
		EmailResponsavel:    paciente.EmailResponsavel,
		Observacoes:         paciente.Observacoes,
		Alergias:            paciente.Alergias,
		Medicacoes:          paciente.Medicacoes,
		CriadoEm:            paciente.CriadoEm,
		AtualizadoEm:        paciente.AtualizadoEm,
		CriadoPor:           paciente.CriadoPor,
		AtualizadoPor:       paciente.AtualizadoPor,
	}
}