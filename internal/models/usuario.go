package models

import (
	"time"

	"gorm.io/gorm"
)

type Usuario struct {
	gorm.Model

	Nome                  string    `json:"nome"`
	Email                 string    `gorm:"unique" json:"email"`
	Senha                 string    `json:"senha"`
	Perfil                string    `json:"perfil"`
	ClinicaID             uint      `json:"clinica_id"`
	SupervisorID          uint      `json:"supervisor_id"`
	Ativo                 bool      `gorm:"default:true" json:"ativo"`
	DataInicioInatividade time.Time `json:"data_inicio_inatividade"`
	DataFimInatividade    time.Time `json:"data_fim_inatividade"`
	MotivoInatividade     string    `json:"motivo_inatividade"`
}

func (Usuario) TableName() string {
	return "usuarios"
}

// CreateUserRequest representa os dados para criar um novo usuário
type CreateUserRequest struct {
	Nome   string `json:"nome" binding:"required,min=3,max=50"`
	Email  string `json:"email" binding:"required,email"`
	Senha  string `json:"senha" binding:"required,min=6"`
	Perfil string `json:"perfil" binding:"required"`
}

// UpdateUserRequest representa os dados para atualizar um usuário
type UpdateUserRequest struct {
	Nome                  *string    `json:"nome" binding:"omitempty,min=3,max=50"`
	Email                 *string    `json:"email" binding:"omitempty,email"`
	Perfil                *string    `json:"perfil"`
	ClinicaID             *uint      `json:"clinica_id"`
	SupervisorID          *uint      `json:"supervisor_id"`
	Ativo                 *bool      `json:"ativo"`
	DataInicioInatividade *time.Time `json:"data_inicio_inatividade"`
	DataFimInatividade    *time.Time `json:"data_fim_inatividade"`
	MotivoInatividade     *string    `json:"motivo_inatividade"`
}

// ChangePasswordRequest representa os dados para alterar a senha
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}
