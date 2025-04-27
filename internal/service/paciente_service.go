// internal/service/paciente_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"clinica_server/internal/models"
	"clinica_server/internal/repository"
	"clinica_server/internal/validator"
)

// Erros específicos do serviço
var (
	ErrInvalidInput      = errors.New("dados de entrada inválidos")
	ErrNotFound          = errors.New("paciente não encontrado")
	ErrUnauthorized      = errors.New("operação não autorizada")
	ErrInternalServer    = errors.New("erro interno do servidor")
	ErrDuplicateResource = errors.New("recurso duplicado")
)

// PacienteService define a interface para operações de serviço de pacientes
type PacienteService interface {
	// Operações CRUD
	Create(ctx context.Context, request models.CreatePacienteRequest) (*models.Paciente, error)
	GetByID(ctx context.Context, id uint) (*models.PacienteDetailDTO, error)
	Update(ctx context.Context, id uint, request models.UpdatePacienteRequest) (*models.Paciente, error)
	Delete(ctx context.Context, id uint) error
	
	// Operações de listagem e busca
	List(ctx context.Context, page, limit int, filters map[string]interface{}) (*models.PacienteListResponse, error)
	Search(ctx context.Context, query string, page, limit int) (*models.PacienteListResponse, error)
	GetByCPF(ctx context.Context, cpf string) (*models.PacienteDetailDTO, error)
	
	// Operações específicas de negócio
	CalcularIdade(ctx context.Context, id uint) (int, error)
	AtualizarDiagnostico(ctx context.Context, id uint, diagnostico string) error
}

// DefaultPacienteService implementa PacienteService
type DefaultPacienteService struct {
	repo      repository.PacienteRepository
	validator validator.PacienteValidator
	converter models.PacienteDTOConverter
}

// NewPacienteService cria uma nova instância de PacienteService
func NewPacienteService(
	repo repository.PacienteRepository,
	validator validator.PacienteValidator,
	converter models.PacienteDTOConverter,
) PacienteService {
	return &DefaultPacienteService{
		repo:      repo,
		validator: validator,
		converter: converter,
	}
}

// Create cria um novo paciente
func (s *DefaultPacienteService) Create(ctx context.Context, request models.CreatePacienteRequest) (*models.Paciente, error) {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	// Converter request para modelo
	paciente := &models.Paciente{
		Nome:                request.Nome,
		DataNascimento:      request.DataNascimento,
		Genero:              request.Genero,
		CPF:                 request.CPF,
		RG:                  request.RG,
		Diagnostico:         request.Diagnostico,
		Telefone:            request.Telefone,
		Email:               request.Email,
		Endereco:            request.Endereco,
		Cidade:              request.Cidade,
		Estado:              request.Estado,
		CEP:                 request.CEP,
		NomeResponsavel:     request.NomeResponsavel,
		TelefoneResponsavel: request.TelefoneResponsavel,
		EmailResponsavel:    request.EmailResponsavel,
		Observacoes:         request.Observacoes,
		Alergias:            request.Alergias,
		Medicacoes:          request.Medicacoes,
	}
	
	// Validar dados (o validador irá definir CriadoPor e AtualizadoPor)
	if err := s.validator.ValidateCreate(ctx, paciente); err != nil {
		// Mapear erros específicos para erros do serviço
		switch {
		case errors.Is(err, validator.ErrCPFInvalid) || errors.Is(err, validator.ErrEmailInvalid):
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		case errors.Is(err, validator.ErrCPFDuplicate):
			return nil, fmt.Errorf("%w: %v", ErrDuplicateResource, err)
		case errors.Is(err, validator.ErrUnauthorized):
			return nil, fmt.Errorf("%w: %v", ErrUnauthorized, err)
		default:
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
	}
	
	// Criar paciente no repositório
	if err := s.repo.Create(ctx, paciente); err != nil {
		// Mapear erros do repositório
		switch {
		case errors.Is(err, repository.ErrDuplicateCPF):
			return nil, fmt.Errorf("%w: CPF já cadastrado", ErrDuplicateResource)
		case errors.Is(err, repository.ErrDatabase):
			return nil, fmt.Errorf("%w: erro ao criar paciente", ErrInternalServer)
		default:
			return nil, fmt.Errorf("%w: %v", ErrInternalServer, err)
		}
	}
	
	return paciente, nil
}

// GetByID busca um paciente pelo ID
func (s *DefaultPacienteService) GetByID(ctx context.Context, id uint) (*models.PacienteDetailDTO, error) {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Validar acesso
	if err := s.validator.ValidateAccess(ctx, id, "view"); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	
	// Buscar paciente no repositório
	paciente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Mapear erros do repositório
		switch {
		case errors.Is(err, repository.ErrPacienteNotFound):
			return nil, ErrNotFound
		case errors.Is(err, repository.ErrInvalidID):
			return nil, fmt.Errorf("%w: ID inválido", ErrInvalidInput)
		default:
			return nil, fmt.Errorf("%w: %v", ErrInternalServer, err)
		}
	}
	
	// Converter para DTO
	detailDTO := s.converter.ToDetailDTO(paciente)
	
	return &detailDTO, nil
}

// Update atualiza um paciente existente
func (s *DefaultPacienteService) Update(ctx context.Context, id uint, request models.UpdatePacienteRequest) (*models.Paciente, error) {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	// Buscar paciente existente
	existingPaciente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Mapear erros do repositório
		switch {
		case errors.Is(err, repository.ErrPacienteNotFound):
			return nil, ErrNotFound
		case errors.Is(err, repository.ErrInvalidID):
			return nil, fmt.Errorf("%w: ID inválido", ErrInvalidInput)
		default:
			return nil, fmt.Errorf("%w: %v", ErrInternalServer, err)
		}
	}
	
	// Atualizar apenas os campos fornecidos
	paciente := &models.Paciente{
		ID:                  id,
		Nome:                existingPaciente.Nome,
		DataNascimento:      existingPaciente.DataNascimento,
		Genero:              existingPaciente.Genero,
		CPF:                 existingPaciente.CPF,
		RG:                  existingPaciente.RG,
		Diagnostico:         existingPaciente.Diagnostico,
		Telefone:            existingPaciente.Telefone,
		Email:               existingPaciente.Email,
		Endereco:            existingPaciente.Endereco,
		Cidade:              existingPaciente.Cidade,
		Estado:              existingPaciente.Estado,
		CEP:                 existingPaciente.CEP,
		NomeResponsavel:     existingPaciente.NomeResponsavel,
		TelefoneResponsavel: existingPaciente.TelefoneResponsavel,
		EmailResponsavel:    existingPaciente.EmailResponsavel,
		Observacoes:         existingPaciente.Observacoes,
		Alergias:            existingPaciente.Alergias,
		Medicacoes:          existingPaciente.Medicacoes,
		CriadoEm:            existingPaciente.CriadoEm,
		CriadoPor:           existingPaciente.CriadoPor,
	}
	
	// Atualizar campos se fornecidos na request
	if request.Nome != "" {
		paciente.Nome = request.Nome
	}
	if !request.DataNascimento.IsZero() {
		paciente.DataNascimento = request.DataNascimento
	}
	if request.Genero != "" {
		paciente.Genero = request.Genero
	}
	if request.CPF != "" {
		paciente.CPF = request.CPF
	}
	if request.RG != "" {
		paciente.RG = request.RG
	}
	if request.Diagnostico != "" {
		paciente.Diagnostico = request.Diagnostico
	}
	if request.Telefone != "" {
		paciente.Telefone = request.Telefone
	}
	if request.Email != "" {
		paciente.Email = request.Email
	}
	if request.Endereco != "" {
		paciente.Endereco = request.Endereco
	}
	if request.Cidade != "" {
		paciente.Cidade = request.Cidade
	}
	if request.Estado != "" {
		paciente.Estado = request.Estado
	}
	if request.CEP != "" {
		paciente.CEP = request.CEP
	}
	if request.NomeResponsavel != "" {
		paciente.NomeResponsavel = request.NomeResponsavel
	}
	if request.TelefoneResponsavel != "" {
		paciente.TelefoneResponsavel = request.TelefoneResponsavel
	}
	if request.EmailResponsavel != "" {
		paciente.EmailResponsavel = request.EmailResponsavel
	}
	if request.Observacoes != "" {
		paciente.Observacoes = request.Observacoes
	}
	if request.Alergias != "" {
		paciente.Alergias = request.Alergias
	}
	if request.Medicacoes != "" {
		paciente.Medicacoes = request.Medicacoes
	}
	
	// Validar dados (o validador irá definir AtualizadoPor)
	if err := s.validator.ValidateUpdate(ctx, paciente); err != nil {
		// Mapear erros específicos para erros do serviço
		switch {
		case errors.Is(err, validator.ErrCPFInvalid) || errors.Is(err, validator.ErrEmailInvalid):
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		case errors.Is(err, validator.ErrCPFDuplicate):
			return nil, fmt.Errorf("%w: %v", ErrDuplicateResource, err)
		case errors.Is(err, validator.ErrUnauthorized):
			return nil, fmt.Errorf("%w: %v", ErrUnauthorized, err)
		default:
			return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
	}
	
	// Atualizar paciente no repositório
	if err := s.repo.Update(ctx, paciente); err != nil {
		// Mapear erros do repositório
		switch {
		case errors.Is(err, repository.ErrPacienteNotFound):
			return nil, ErrNotFound
		case errors.Is(err, repository.ErrDuplicateCPF):
			return nil, fmt.Errorf("%w: CPF já cadastrado", ErrDuplicateResource)
		case errors.Is(err, repository.ErrDatabase):
			return nil, fmt.Errorf("%w: erro ao atualizar paciente", ErrInternalServer)
		default:
			return nil, fmt.Errorf("%w: %v", ErrInternalServer, err)
		}
	}
	
	return paciente, nil
}

// Delete remove um paciente
func (s *DefaultPacienteService) Delete(ctx context.Context, id uint) error {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Validar exclusão
	if err := s.validator.ValidateDelete(ctx, id); err != nil {
		// Mapear erros específicos para erros do serviço
		switch {
		case errors.Is(err, validator.ErrUnauthorized):
			return fmt.Errorf("%w: %v", ErrUnauthorized, err)
		default:
			return fmt.Errorf("%w: %v", ErrInvalidInput, err)
		}
	}
	
	// Obter usuário do contexto
	user, err := s.validator.GetUserFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	
	// Excluir paciente no repositório
	if err := s.repo.Delete(ctx, id, user.UserID); err != nil {
		// Mapear erros do repositório
		switch {
		case errors.Is(err, repository.ErrPacienteNotFound):
			return ErrNotFound
		case errors.Is(err, repository.ErrInvalidID):
			return fmt.Errorf("%w: ID inválido", ErrInvalidInput)
		default:
			return fmt.Errorf("%w: %v", ErrInternalServer, err)
		}
	}
	
	return nil
}

// List retorna uma lista paginada de pacientes
func (s *DefaultPacienteService) List(ctx context.Context, page, limit int, filters map[string]interface{}) (*models.PacienteListResponse, error) {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	// Validar acesso
	if err := s.validator.ValidateAccess(ctx, 0, "view"); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	
	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	// Buscar pacientes no repositório
	pacientes, total, err := s.repo.List(ctx, page, limit, filters)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	
	// Converter para resposta paginada
	response := s.converter.ToPacienteListResponse(pacientes, page, limit, total)
	
	return &response, nil
}

// Search busca pacientes por um termo de busca
func (s *DefaultPacienteService) Search(ctx context.Context, query string, page, limit int) (*models.PacienteListResponse, error) {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	// Validar acesso
	if err := s.validator.ValidateAccess(ctx, 0, "view"); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	
	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	// Buscar pacientes no repositório
	pacientes, total, err := s.repo.Search(ctx, query, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	
	// Converter para resposta paginada
	response := s.converter.ToPacienteListResponse(pacientes, page, limit, total)
	
	return &response, nil
}

// GetByCPF busca um paciente pelo CPF
func (s *DefaultPacienteService) GetByCPF(ctx context.Context, cpf string) (*models.PacienteDetailDTO, error) {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Validar acesso
	if err := s.validator.ValidateAccess(ctx, 0, "view"); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	
	// Validar CPF
	if err := s.validator.ValidateCPF(cpf); err != nil {
		return nil, fmt.Errorf("%w: CPF inválido", ErrInvalidInput)
	}
	
	// Buscar paciente no repositório
	paciente, err := s.repo.GetByCPF(ctx, cpf)
	if err != nil {
		// Mapear erros do repositório
		switch {
		case errors.Is(err, repository.ErrPacienteNotFound):
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("%w: %v", ErrInternalServer, err)
		}
	}
	
	// Converter para DTO
	detailDTO := s.converter.ToDetailDTO(paciente)
	
	return &detailDTO, nil
}

// CalcularIdade calcula a idade de um paciente
func (s *DefaultPacienteService) CalcularIdade(ctx context.Context, id uint) (int, error) {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Validar acesso
	if err := s.validator.ValidateAccess(ctx, id, "view"); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	
	// Buscar paciente no repositório
	paciente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Mapear erros do repositório
		switch {
		case errors.Is(err, repository.ErrPacienteNotFound):
			return 0, ErrNotFound
		case errors.Is(err, repository.ErrInvalidID):
			return 0, fmt.Errorf("%w: ID inválido", ErrInvalidInput)
		default:
			return 0, fmt.Errorf("%w: %v", ErrInternalServer, err)
		}
	}
	
	// Calcular idade usando a função do pacote validacao
	idade := s.converter.ToDTO(paciente).Idade
	
	return idade, nil
}

// AtualizarDiagnostico atualiza apenas o diagnóstico de um paciente
func (s *DefaultPacienteService) AtualizarDiagnostico(ctx context.Context, id uint, diagnostico string) error {
	// Criar timeout para a operação
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Validar acesso (poderia ter uma permissão específica para atualizar diagnóstico)
	if err := s.validator.ValidateAccess(ctx, id, "update"); err != nil {
		return fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	
	// Buscar paciente existente
	existingPaciente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Mapear erros do repositório
		switch {
		case errors.Is(err, repository.ErrPacienteNotFound):
			return ErrNotFound
		case errors.Is(err, repository.ErrInvalidID):
			return fmt.Errorf("%w: ID inválido", ErrInvalidInput)
		default:
			return fmt.Errorf("%w: %v", ErrInternalServer, err)
		}
	}
	
	// Obter usuário do contexto
	user, err := s.validator.GetUserFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnauthorized, err)
	}
	
	// Atualizar apenas o diagnóstico
	existingPaciente.Diagnostico = diagnostico
	existingPaciente.AtualizadoPor = user.UserID
	
	// Atualizar paciente no repositório
	if err := s.repo.Update(ctx, existingPaciente); err != nil {
		return fmt.Errorf("%w: %v", ErrInternalServer, err)
	}
	
	return nil
}