// internal/models/paciente_dto.go
package models

import (
	"clinica_server/internal/utils/validacao"
	"time"
)

// PacienteDTO representa a versão simplificada do paciente para listagens
type PacienteDTO struct {
	ID             uint      `json:"id"`
	Nome           string    `json:"nome"`
	DataNascimento time.Time `json:"data_nascimento"`
	Idade          int       `json:"idade"`
	Genero         string    `json:"genero"`
	CPF            string    `json:"cpf"`
	Telefone       string    `json:"telefone"`
	Diagnostico    string    `json:"diagnostico"`
	CriadoEm       time.Time `json:"criado_em"`
	CriadoPor      uint      `json:"criado_por"`
}

// PacienteDetailDTO representa a versão detalhada do paciente
type PacienteDetailDTO struct {
	ID                 uint      `json:"id"`
	Nome               string    `json:"nome"`
	DataNascimento     time.Time `json:"data_nascimento"`
	Idade              int       `json:"idade"`
	Genero             string    `json:"genero"`
	CPF                string    `json:"cpf"`
	RG                 string    `json:"rg"`
	Diagnostico        string    `json:"diagnostico"`
	Telefone           string    `json:"telefone"`
	Email              string    `json:"email"`
	Endereco           string    `json:"endereco"`
	Cidade             string    `json:"cidade"`
	Estado             string    `json:"estado"`
	CEP                string    `json:"cep"`
	NomeResponsavel    string    `json:"nome_responsavel"`
	TelefoneResponsavel string   `json:"telefone_responsavel"`
	EmailResponsavel   string    `json:"email_responsavel"`
	Observacoes        string    `json:"observacoes"`
	Alergias           string    `json:"alergias"`
	Medicacoes         string    `json:"medicacoes"`
	CriadoEm           time.Time `json:"criado_em"`
	CriadoPor          uint      `json:"criado_por"`
	AtualizadoEm       time.Time `json:"atualizado_em"`
	AtualizadoPor      uint      `json:"atualizado_por"`
}

// PacienteListResponse representa a resposta paginada de pacientes
type PacienteListResponse struct {
	Data  []PacienteDTO `json:"data"`
	Meta  PaginationMeta `json:"meta"`
}

// PaginationMeta contém informações sobre a paginação
type PaginationMeta struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int   `json:"total_pages"`
}

// PacienteDTOConverter define a interface para conversão entre Paciente e DTOs
type PacienteDTOConverter interface {
	ToDTO(paciente *Paciente) PacienteDTO
	ToDetailDTO(paciente *Paciente) PacienteDetailDTO
	ToDTOList(pacientes []Paciente) []PacienteDTO
	ToPacienteListResponse(pacientes []Paciente, page, perPage int, total int64) PacienteListResponse
}

// DefaultPacienteDTOConverter implementa PacienteDTOConverter
type DefaultPacienteDTOConverter struct{}

// NewPacienteDTOConverter cria uma nova instância do conversor de DTOs
func NewPacienteDTOConverter() PacienteDTOConverter {
	return &DefaultPacienteDTOConverter{}
}

// ToDTO converte um Paciente para PacienteDTO
func (c *DefaultPacienteDTOConverter) ToDTO(paciente *Paciente) PacienteDTO {
	if paciente == nil {
		return PacienteDTO{}
	}

	return PacienteDTO{
		ID:             paciente.ID,
		Nome:           paciente.Nome,
		DataNascimento: paciente.DataNascimento,
		Idade:          validacao.Calcularidade(paciente.DataNascimento),
		Genero:         paciente.Genero,
		CPF:            paciente.CPF,
		Telefone:       paciente.Telefone,
		Diagnostico:    paciente.Diagnostico,
		CriadoEm:       paciente.CriadoEm,
		CriadoPor:      paciente.CriadoPor,
	}
}

// ToDetailDTO converte um Paciente para PacienteDetailDTO
func (c *DefaultPacienteDTOConverter) ToDetailDTO(paciente *Paciente) PacienteDetailDTO {
	if paciente == nil {
		return PacienteDetailDTO{}
	}

	return PacienteDetailDTO{
		ID:                 paciente.ID,
		Nome:               paciente.Nome,
		DataNascimento:     paciente.DataNascimento,
		Idade:              validacao.Calcularidade(paciente.DataNascimento),
		Genero:             paciente.Genero,
		CPF:                paciente.CPF,
		RG:                 paciente.RG,
		Diagnostico:        paciente.Diagnostico,
		Telefone:           paciente.Telefone,
		Email:              paciente.Email,
		Endereco:           paciente.Endereco,
		Cidade:             paciente.Cidade,
		Estado:             paciente.Estado,
		CEP:                paciente.CEP,
		NomeResponsavel:    paciente.NomeResponsavel,
		TelefoneResponsavel: paciente.TelefoneResponsavel,
		EmailResponsavel:   paciente.EmailResponsavel,
		Observacoes:        paciente.Observacoes,
		Alergias:           paciente.Alergias,
		Medicacoes:         paciente.Medicacoes,
		CriadoEm:           paciente.CriadoEm,
		CriadoPor:          paciente.CriadoPor,
		AtualizadoEm:       paciente.AtualizadoEm,
		AtualizadoPor:      paciente.AtualizadoPor,
	}
}

// ToDTOList converte uma lista de Paciente para uma lista de PacienteDTO
func (c *DefaultPacienteDTOConverter) ToDTOList(pacientes []Paciente) []PacienteDTO {
	dtos := make([]PacienteDTO, len(pacientes))
	for i, paciente := range pacientes {
		pacienteCopy := paciente // Cria uma cópia para evitar problemas com o loop
		dtos[i] = c.ToDTO(&pacienteCopy)
	}
	return dtos
}

// ToPacienteListResponse cria uma resposta paginada de pacientes
func (c *DefaultPacienteDTOConverter) ToPacienteListResponse(pacientes []Paciente, page, perPage int, total int64) PacienteListResponse {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return PacienteListResponse{
		Data: c.ToDTOList(pacientes),
		Meta: PaginationMeta{
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		},
	}
}