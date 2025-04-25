package models

import (
	"time"
	"msd-services clinica_server/internal/utils/validacao"

	"gorm.io/gorm"
)

type Paciente struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Nome               string         `gorm:"not null;size:100" json:"nome"`
	DataNascimento     time.Time      `json:"data_nascimento"`
	Genero             string         `gorm:"size:20" json:"genero"`
	CPF                string         `gorm:"size:14;uniqueIndex" json:"cpf"`
	RG                 string         `gorm:"size:20" json:"rg"`
	Diagnostico        string         `gorm:"size:255" json:"diagnostico"`
	Telefone           string         `gorm:"size:20" json:"telefone"`
	Email              string         `gorm:"size:100" json:"email"`
	Endereco           string         `gorm:"size:255" json:"endereco"`
	Cidade             string         `gorm:"size:100" json:"cidade"`
	Estado             string         `gorm:"size:2" json:"estado"`
	CEP                string         `gorm:"size:10" json:"cep"`
	NomeResponsavel    string         `gorm:"size:100" json:"nome_responsavel"`
	TelefoneResponsavel string        `gorm:"size:20" json:"telefone_responsavel"`
	EmailResponsavel   string         `gorm:"size:100" json:"email_responsavel"`
	Observacoes        string         `gorm:"type:text" json:"observacoes"`
	Alergias           string         `gorm:"type:text" json:"alergias"`
	Medicacoes         string         `gorm:"type:text" json:"medicacoes"`
	CriadoEm           time.Time      `gorm:"autoCreateTime" json:"criado_em"`
	AtualizadoEm       time.Time      `gorm:"autoUpdateTime" json:"atualizado_em"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

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
	NomeResponsavel     string    `json:"nome_responsavel" binding:"required,max=100"`
	TelefoneResponsavel string    `json:"telefone_responsavel" binding:"required,max=20"`
	EmailResponsavel    string    `json:"email_responsavel" binding:"omitempty,email,max=100"`
	Observacoes         string    `json:"observacoes"`
	Alergias            string    `json:"alergias"`
	Medicacoes          string    `json:"medicacoes"`
}

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

// BeforeSave é um hook do GORM que é executado antes de salvar o registro
func (p *Paciente) BeforeSave(tx *gorm.DB) error {
	// Validação e formatação do CPF
	if p.CPF != "" {
		// Formata o CPF
		p.CPF = validators.FormatarCPF(p.CPF)
		
		// Valida o CPF
		if !validators.ValidarCPF(p.CPF) {
			return tx.AddError(gorm.ErrInvalidData).AddError(
				&gorm.ErrInvalidField{Field: "CPF", Message: "CPF inválido"})
		}
	}
	
	// Validação de email (opcional, pode ser feita no pacote validators também)
	if p.Email != "" {
		// Aqui você poderia usar uma função validators.ValidarEmail(p.Email)
	}
	
	return nil
}

// BeforeCreate é um hook do GORM que é executado antes de criar um novo registro
func (p *Paciente) BeforeCreate(tx *gorm.DB) error {
	// Validações específicas para criação
	if p.Nome == "" {
		return tx.AddError(gorm.ErrInvalidData).AddError(
			&gorm.ErrInvalidField{Field: "Nome", Message: "Nome é obrigatório"})
	}
	
	// Outras validações específicas para criação
	
	return nil
}

// AfterFind é um hook do GORM que é executado após buscar o registro do banco
func (p *Paciente) AfterFind(tx *gorm.DB) error {
	// Você pode realizar operações após carregar o registro
	// Por exemplo, formatar dados para exibição
	return nil
}

// CalcularIdade calcula a idade do paciente com base na data de nascimento
func (p *Paciente) CalcularIdade() int {
	if p.DataNascimento.IsZero() {
		return 0
	}
	
	now := time.Now()
	idade := now.Year() - p.DataNascimento.Year()
	
	// Ajusta a idade se ainda não fez aniversário este ano
	if now.Month() < p.DataNascimento.Month() || 
	   (now.Month() == p.DataNascimento.Month() && now.Day() < p.DataNascimento.Day()) {
		idade--
	}
	
	return idade
}

// TableName especifica o nome da tabela no banco de dados
func (Paciente) TableName() string {
	return "pacientes"
}