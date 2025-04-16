package models

import "time"

type Paciente struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Nome           string    `gorm:"not null" json:"nome"`
	DataNascimento time.Time `json:"data_nascimento"`
	Genero         string    `json:"genero"`
	Diagnostico    string    `json:"diagnostico"`
	Telefone       string    `json:"telefone"`
	Email          string    `json:"email"`
	Endereco       string    `json:"endereco"`
	CriadoEm       time.Time `gorm:"autoCreateTime" json:"criado_em"`
}

type CreatePacienteRequest struct {
	Nome           string    `json:"nome" binding:"required"`
	DataNascimento time.Time `json:"data_nascimento"`
	Genero         string    `json:"genero"`
	Diagnostico    string    `json:"diagnostico"`
	Telefone       string    `json:"telefone"`
	Email          string    `json:"email"`
	Endereco       string    `json:"endereco"`
}

type UpdatePacienteRequest struct {
	Nome           string    `json:"nome"`
	DataNascimento time.Time `json:"data_nascimento"`
	Genero         string    `json:"genero"`
	Diagnostico    string    `json:"diagnostico"`
	Telefone       string    `json:"telefone"`
	Email          string    `json:"email"`
	Endereco       string    `json:"endereco"`
}