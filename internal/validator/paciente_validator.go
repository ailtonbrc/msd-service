package validator

import (
	"context"
	"errors"
	"fmt"

	"clinica_server/internal/auth"
	"clinica_server/internal/models"
	"clinica_server/internal/repository"
	"clinica_server/internal/utils/validacao"
)

// PacienteValidator define a interface para validação de pacientes
type PacienteValidator interface {
	// ValidateCreate valida os dados para criação de um paciente
	ValidateCreate(ctx context.Context, paciente *models.Paciente) error
	
	// ValidateUpdate valida os dados para atualização de um paciente
	ValidateUpdate(ctx context.Context, paciente *models.Paciente) error
	
	// ValidateDelete valida a exclusão de um paciente
	ValidateDelete(ctx context.Context, id uint) error
	
	// ValidateAccess verifica se o usuário tem permissão para acessar o paciente
	ValidateAccess(ctx context.Context, id uint, action string) error
}

// DefaultPacienteValidator implementa PacienteValidator
type DefaultPacienteValidator struct {
	pacienteRepo repository.PacienteRepository
}

// NewPacienteValidator cria uma nova instância de PacienteValidator
func NewPacienteValidator(pacienteRepo repository.PacienteRepository) PacienteValidator {
	return &DefaultPacienteValidator{
		pacienteRepo: pacienteRepo,
	}
}

// ValidateCreate valida os dados para criação de um paciente
// Recebe o contexto e o modelo do paciente a ser validado
// Retorna erro se a validação falhar
func (v *DefaultPacienteValidator) ValidateCreate(ctx context.Context, paciente *models.Paciente) error {
	// Verificar permissão
	if err := v.ValidateAccess(ctx, 0, "create"); err != nil {
		return err
	}

	// Validar campos obrigatórios
	if paciente.Nome == "" {
		return errors.New("nome é obrigatório")
	}

	if paciente.Telefone == "" {
		return errors.New("telefone é obrigatório")
	}

	// Validar CPF se fornecido
	if paciente.CPF != "" {
		if !validacao.ValidarCPF(paciente.CPF) {
			return errors.New("CPF inválido")
		}

		// Verificar se já existe um paciente com o mesmo CPF
		exists, err := v.pacienteRepo.ExistsByCPF(ctx, paciente.CPF, 0)
		if err != nil {
			return fmt.Errorf("erro ao verificar CPF: %w", err)
		}
		if exists {
			return errors.New("já existe um paciente com este CPF")
		}
	}

	// Validar email se fornecido
	if paciente.Email != "" && !validacao.ValidarEmail(paciente.Email) {
		return errors.New("email inválido")
	}

	// Validar CEP se fornecido
	if paciente.CEP != "" && !validacao.ValidarCEP(paciente.CEP) {
		return errors.New("CEP inválido")
	}

	// Validar telefone
	if !validacao.ValidarTelefone(paciente.Telefone) {
		return errors.New("telefone inválido")
	}

	// Validar RG se fornecido
	if paciente.RG != "" && !validacao.ValidarRG(paciente.RG) {
		return errors.New("RG inválido")
	}

	return nil
}

// ValidateUpdate valida os dados para atualização de um paciente
// Recebe o contexto e o modelo do paciente a ser validado
// Retorna erro se a validação falhar
func (v *DefaultPacienteValidator) ValidateUpdate(ctx context.Context, paciente *models.Paciente) error {
	// Verificar permissão
	if err := v.ValidateAccess(ctx, paciente.ID, "update"); err != nil {
		return err
	}

	// Validar CPF se fornecido
	if paciente.CPF != "" {
		if !validacao.ValidarCPF(paciente.CPF) {
			return errors.New("CPF inválido")
		}

		// Verificar se já existe outro paciente com o mesmo CPF
		exists, err := v.pacienteRepo.ExistsByCPF(ctx, paciente.CPF, paciente.ID)
		if err != nil {
			return fmt.Errorf("erro ao verificar CPF: %w", err)
		}
		if exists {
			return errors.New("já existe outro paciente com este CPF")
		}
	}

	// Validar email se fornecido
	if paciente.Email != "" && !validacao.ValidarEmail(paciente.Email) {
		return errors.New("email inválido")
	}

	// Validar CEP se fornecido
	if paciente.CEP != "" && !validacao.ValidarCEP(paciente.CEP) {
		return errors.New("CEP inválido")
	}

	// Validar telefone se fornecido
	if paciente.Telefone != "" && !validacao.ValidarTelefone(paciente.Telefone) {
		return errors.New("telefone inválido")
	}

	// Validar RG se fornecido
	if paciente.RG != "" && !validacao.ValidarRG(paciente.RG) {
		return errors.New("RG inválido")
	}

	return nil
}

// ValidateDelete valida a exclusão de um paciente
// Recebe o contexto e o ID do paciente a ser excluído
// Retorna erro se a validação falhar
func (v *DefaultPacienteValidator) ValidateDelete(ctx context.Context, id uint) error {
	// Verificar permissão
	if err := v.ValidateAccess(ctx, id, "delete"); err != nil {
		return err
	}

	// Verificar se o paciente existe
	paciente, err := v.pacienteRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao buscar paciente: %w", err)
	}
	if paciente == nil {
		return errors.New("paciente não encontrado")
	}

	return nil
}

// ValidateAccess verifica se o usuário tem permissão para acessar o paciente
// Recebe o contexto, o ID do paciente e a ação a ser realizada
// Retorna erro se o usuário não tiver permissão
func (v *DefaultPacienteValidator) ValidateAccess(ctx context.Context, id uint, action string) error {
	// Obter usuário do contexto
	userClaims, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return errors.New("usuário não autenticado")
	}


	
	// Verificar permissão
	permission := fmt.Sprintf("pacientes:%s", action)
	if !auth.HasPermission(userClaims, permission) {
		return fmt.Errorf("sem permissão para %s pacientes", action)
	}

	return nil
}