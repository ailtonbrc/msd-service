// internal/validator/paciente_validator.go
package validator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"clinica_server/internal/auth"
	"clinica_server/internal/models"
	"clinica_server/internal/repository"
	"clinica_server/utils/validacao"
)

// Erros de validação específicos para pacientes
var (
	ErrNomeRequired           = errors.New("nome é obrigatório")
	ErrDataNascimentoRequired = errors.New("data de nascimento é obrigatória")
	ErrDataNascimentoInvalid  = errors.New("data de nascimento não pode ser no futuro")
	ErrIdadeMinima            = errors.New("paciente deve ter pelo menos 1 ano de idade")
	ErrCPFInvalid             = errors.New("CPF inválido")
	ErrCPFDuplicate           = errors.New("CPF já cadastrado para outro paciente")
	ErrEmailInvalid           = errors.New("formato de email inválido")
	ErrTelefoneRequired       = errors.New("telefone é obrigatório")
	ErrResponsavelRequired    = errors.New("nome do responsável é obrigatório para pacientes menores de idade")
	ErrTelResponsavelRequired = errors.New("telefone do responsável é obrigatório para pacientes menores de idade")
	ErrUnauthorized           = errors.New("usuário não autorizado para esta operação")
	ErrUserNotFound           = errors.New("usuário não encontrado no contexto")
)

// PacienteValidator define a interface para validação de pacientes
type PacienteValidator interface {
	// Validações para criação e atualização
	ValidateCreate(ctx context.Context, paciente *models.Paciente) error
	ValidateUpdate(ctx context.Context, paciente *models.Paciente) error
	
	// Validações para operações específicas
	ValidateDelete(ctx context.Context, id uint) error
	ValidateAccess(ctx context.Context, pacienteID uint, action string) error
	
	// Validações de campos específicos
	ValidateCPF(cpf string) error
	ValidateEmail(email string) error
	ValidateIdade(dataNascimento time.Time) error
	
	// Obter usuário do contexto
	GetUserFromContext(ctx context.Context) (*auth.UserClaims, error)
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
func (v *DefaultPacienteValidator) ValidateCreate(ctx context.Context, paciente *models.Paciente) error {
	// Verificar permissão do usuário (usando JWT)
	if err := v.ValidateAccess(ctx, 0, "create"); err != nil {
		return err
	}
	
	// Obter usuário do contexto
	user, err := v.GetUserFromContext(ctx)
	if err != nil {
		return err
	}
	
	// Definir o usuário que está criando o paciente
	paciente.CriadoPor = user.UserID
	paciente.AtualizadoPor = user.UserID
	
	// Validar campos obrigatórios
	if paciente.Nome == "" {
		return ErrNomeRequired
	}
	
	if paciente.DataNascimento.IsZero() {
		return ErrDataNascimentoRequired
	}
	
	// Validar data de nascimento
	if err := v.ValidateIdade(paciente.DataNascimento); err != nil {
		return err
	}
	
	// Validar CPF
	if paciente.CPF != "" {
		if err := v.ValidateCPF(paciente.CPF); err != nil {
			return err
		}
		
		// Verificar se o CPF já está cadastrado
		exists, err := v.pacienteRepo.ExistsByCPF(ctx, paciente.CPF, 0)
		if err != nil {
			return fmt.Errorf("erro ao verificar CPF: %w", err)
		}
		if exists {
			return ErrCPFDuplicate
		}
	}
	
	// Validar email
	if paciente.Email != "" {
		if err := v.ValidateEmail(paciente.Email); err != nil {
			return err
		}
	}
	
	// Validar telefone
	if paciente.Telefone == "" {
		return ErrTelefoneRequired
	}
	
	// Validar dados do responsável para menores de idade
	idade := validacao.CalculateAge(paciente.DataNascimento)
	if idade < 18 {
		if paciente.NomeResponsavel == "" {
			return ErrResponsavelRequired
		}
		
		if paciente.TelefoneResponsavel == "" {
			return ErrTelResponsavelRequired
		}
	}
	
	return nil
}

// ValidateUpdate valida os dados para atualização de um paciente
func (v *DefaultPacienteValidator) ValidateUpdate(ctx context.Context, paciente *models.Paciente) error {
	// Verificar permissão do usuário (usando JWT)
	if err := v.ValidateAccess(ctx, paciente.ID, "update"); err != nil {
		return err
	}
	
	// Obter usuário do contexto
	user, err := v.GetUserFromContext(ctx)
	if err != nil {
		return err
	}
	
	// Definir o usuário que está atualizando o paciente
	paciente.AtualizadoPor = user.UserID
	
	// Verificar se o paciente existe
	existingPaciente, err := v.pacienteRepo.GetByID(ctx, paciente.ID)
	if err != nil {
		return err
	}
	
	// Para campos que não foram fornecidos na atualização, usar os valores existentes
	if paciente.Nome == "" {
		paciente.Nome = existingPaciente.Nome
	}
	
	if paciente.DataNascimento.IsZero() {
		paciente.DataNascimento = existingPaciente.DataNascimento
	} else {
		// Validar data de nascimento
		if err := v.ValidateIdade(paciente.DataNascimento); err != nil {
			return err
		}
	}
	
	// Validar CPF se foi fornecido
	if paciente.CPF != "" && paciente.CPF != existingPaciente.CPF {
		if err := v.ValidateCPF(paciente.CPF); err != nil {
			return err
		}
		
		// Verificar se o CPF já está cadastrado para outro paciente
		exists, err := v.pacienteRepo.ExistsByCPF(ctx, paciente.CPF, paciente.ID)
		if err != nil {
			return fmt.Errorf("erro ao verificar CPF: %w", err)
		}
		if exists {
			return ErrCPFDuplicate
		}
	}
	
	// Validar email se foi fornecido
	if paciente.Email != "" && paciente.Email != existingPaciente.Email {
		if err := v.ValidateEmail(paciente.Email); err != nil {
			return err
		}
	}
	
	// Validar dados do responsável para menores de idade
	idade := validacao.CalculateAge(paciente.DataNascimento)
	if idade < 18 {
		if paciente.NomeResponsavel == "" {
			paciente.NomeResponsavel = existingPaciente.NomeResponsavel
			if paciente.NomeResponsavel == "" {
				return ErrResponsavelRequired
			}
		}
		
		if paciente.TelefoneResponsavel == "" {
			paciente.TelefoneResponsavel = existingPaciente.TelefoneResponsavel
			if paciente.TelefoneResponsavel == "" {
				return ErrTelResponsavelRequired
			}
		}
	}
	
	// Manter o usuário que criou o paciente
	paciente.CriadoPor = existingPaciente.CriadoPor
	
	return nil
}

// ValidateDelete valida se um paciente pode ser excluído
func (v *DefaultPacienteValidator) ValidateDelete(ctx context.Context, id uint) error {
	// Verificar permissão do usuário (usando JWT)
	if err := v.ValidateAccess(ctx, id, "delete"); err != nil {
		return err
	}
	
	// Verificar se o paciente existe
	_, err := v.pacienteRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Aqui poderiam ser adicionadas outras validações de negócio
	// Por exemplo, verificar se o paciente não tem consultas agendadas
	
	return nil
}

// ValidateAccess valida se o usuário tem permissão para a operação
func (v *DefaultPacienteValidator) ValidateAccess(ctx context.Context, pacienteID uint, action string) error {
	// Obter informações do usuário do contexto (JWT)
	userClaims, err := v.GetUserFromContext(ctx)
	if err != nil {
		return err
	}
	
	// Verificar se o usuário tem a permissão necessária
	permissionRequired := fmt.Sprintf("pacientes:%s", action)
	if !auth.HasPermission(userClaims, permissionRequired) {
		return ErrUnauthorized
	}
	
	return nil
}

// ValidateCPF valida o formato e a validade do CPF
func (v *DefaultPacienteValidator) ValidateCPF(cpf string) error {
	if cpf == "" {
		return nil // CPF é opcional
	}
	
	// Usar a função de validação de CPF do pacote validacao
	if !validacao.ValidarCPF(cpf) {
		return ErrCPFInvalid
	}
	
	return nil
}

// ValidateEmail valida o formato do email
func (v *DefaultPacienteValidator) ValidateEmail(email string) error {
	if email == "" {
		return nil // Email é opcional
	}
	
	// Usar uma função de validação de email do pacote validacao
	if !validacao.ValidarEmail(email) {
		return ErrEmailInvalid
	}
	
	return nil
}

// ValidateIdade valida a data de nascimento e a idade
func (v *DefaultPacienteValidator) ValidateIdade(dataNascimento time.Time) error {
	if dataNascimento.IsZero() {
		return ErrDataNascimentoRequired
	}
	
	// Verificar se a data de nascimento não é no futuro
	if dataNascimento.After(time.Now()) {
		return ErrDataNascimentoInvalid
	}
	
	// Verificar idade mínima (exemplo: 1 ano)
	//idade := validacao.CalculateAge(dataNascimento)
	//if idade < 1 {
	//	return ErrIdadeMinima
	//}
	
	return nil
}

// GetUserFromContext extrai os claims do usuário do contexto
func (v *DefaultPacienteValidator) GetUserFromContext(ctx context.Context) (*auth.UserClaims, error) {
	// Extrair os claims do usuário do contexto
	userClaims, ok := auth.GetUserFromContext(ctx)
	if !ok || userClaims == nil {
		return nil, ErrUserNotFound
	}
	
	return userClaims, nil
}