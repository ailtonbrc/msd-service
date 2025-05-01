package service

import (
	"context"
	"errors"

	"clinica_server/internal/auth"
	"clinica_server/internal/models"
	"clinica_server/internal/repository"
	"clinica_server/internal/validator"
)

// PacienteService define a interface para operações de serviço de pacientes
type PacienteService interface {
	// Create cria um novo paciente
	Create(ctx context.Context, req *models.CreatePacienteRequest) (*models.PacienteResponse, error)
	
	// GetByID busca um paciente pelo ID
	GetByID(ctx context.Context, id uint) (*models.PacienteResponse, error)
	
	// GetAll retorna todos os pacientes com paginação
	GetAll(ctx context.Context, page, pageSize int, search string) (*models.PaginatedResponse, error)
	
	// Update atualiza um paciente existente
	Update(ctx context.Context, id uint, req *models.UpdatePacienteRequest) (*models.PacienteResponse, error)
	
	// Delete exclui um paciente pelo ID
	Delete(ctx context.Context, id uint) error
	
	// Search busca pacientes com filtros avançados
	Search(ctx context.Context, filtro *models.PacienteFiltro) (*models.PaginatedResponse, error)
}

// DefaultPacienteService implementa PacienteService
type DefaultPacienteService struct {
	pacienteRepo      repository.PacienteRepository
	pacienteValidator validator.PacienteValidator
	dtoConverter      models.PacienteDTOConverter
}

// NewPacienteService cria uma nova instância de PacienteService
func NewPacienteService(
	pacienteRepo repository.PacienteRepository,
	pacienteValidator validator.PacienteValidator,
	dtoConverter models.PacienteDTOConverter,
) PacienteService {
	return &DefaultPacienteService{
		pacienteRepo:      pacienteRepo,
		pacienteValidator: pacienteValidator,
		dtoConverter:      dtoConverter,
	}
}

// Create cria um novo paciente
// Recebe o contexto e os dados do paciente a ser criado
// Retorna os dados do paciente criado ou erro em caso de falha
func (s *DefaultPacienteService) Create(ctx context.Context, req *models.CreatePacienteRequest) (*models.PacienteResponse, error) {
	// Obter ID do usuário do contexto
	userClaims, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, errors.New("usuário não autenticado")
	}
	
	// Converter DTO para modelo
	paciente := s.dtoConverter.CreateRequestToModel(req, userClaims.UserID)
	
	// Validar dados
	if err := s.pacienteValidator.ValidateCreate(ctx, paciente); err != nil {
		return nil, err
	}
	
	// Salvar no banco de dados
	if err := s.pacienteRepo.Create(ctx, paciente); err != nil {
		return nil, err
	}
	
	// Converter modelo para resposta
	response := s.dtoConverter.ModelToResponse(paciente)
	
	return response, nil
}

// GetByID busca um paciente pelo ID
// Recebe o contexto e o ID do paciente
// Retorna os dados do paciente ou erro em caso de falha
func (s *DefaultPacienteService) GetByID(ctx context.Context, id uint) (*models.PacienteResponse, error) {
	// Verificar permissão
	if err := s.pacienteValidator.ValidateAccess(ctx, id, "read"); err != nil {
		return nil, err
	}
	
	// Buscar paciente
	paciente, err := s.pacienteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if paciente == nil {
		return nil, errors.New("paciente não encontrado")
	}
	
	// Converter modelo para resposta
	response := s.dtoConverter.ModelToResponse(paciente)
	
	return response, nil
}

// GetAll retorna todos os pacientes com paginação
// Recebe o contexto, página, tamanho da página e termo de busca
// Retorna a lista paginada de pacientes ou erro em caso de falha
func (s *DefaultPacienteService) GetAll(ctx context.Context, page, pageSize int, search string) (*models.PaginatedResponse, error) {
	// Verificar permissão
	if err := s.pacienteValidator.ValidateAccess(ctx, 0, "read"); err != nil {
		return nil, err
	}
	
	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	// Buscar pacientes
	pacientes, total, err := s.pacienteRepo.GetAll(ctx, page, pageSize, search)
	if err != nil {
		return nil, err
	}
	
	// Converter modelos para respostas
	var responses []models.PacienteResponse
	for i := range pacientes {
		responses = append(responses, *s.dtoConverter.ModelToResponse(&pacientes[i]))
	}
	
	// Criar resposta paginada
	response := &models.PaginatedResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Data:     responses,
	}
	
	return response, nil
}

// Update atualiza um paciente existente
// Recebe o contexto, o ID do paciente e os dados a serem atualizados
// Retorna os dados do paciente atualizado ou erro em caso de falha
func (s *DefaultPacienteService) Update(ctx context.Context, id uint, req *models.UpdatePacienteRequest) (*models.PacienteResponse, error) {
	// Obter ID do usuário do contexto
	userClaims, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, errors.New("usuário não autenticado")
	}
	
	// Buscar paciente existente
	existingPaciente, err := s.pacienteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if existingPaciente == nil {
		return nil, errors.New("paciente não encontrado")
	}
	
	// Converter DTO para modelo
	paciente := s.dtoConverter.UpdateRequestToModel(req, existingPaciente, userClaims.UserID)
	
	// Validar dados
	if err := s.pacienteValidator.ValidateUpdate(ctx, paciente); err != nil {
		return nil, err
	}
	
	// Atualizar no banco de dados
	if err := s.pacienteRepo.Update(ctx, paciente); err != nil {
		return nil, err
	}
	
	// Converter modelo para resposta
	response := s.dtoConverter.ModelToResponse(paciente)
	
	return response, nil
}

// Delete exclui um paciente pelo ID
// Recebe o contexto e o ID do paciente a ser excluído
// Retorna erro em caso de falha
func (s *DefaultPacienteService) Delete(ctx context.Context, id uint) error {
	// Validar exclusão
	if err := s.pacienteValidator.ValidateDelete(ctx, id); err != nil {
		return err
	}
	
	// Excluir paciente
	return s.pacienteRepo.Delete(ctx, id)
}

// Search busca pacientes com filtros avançados
// Recebe o contexto e os filtros de busca
// Retorna a lista paginada de pacientes ou erro em caso de falha
func (s *DefaultPacienteService) Search(ctx context.Context, filtro *models.PacienteFiltro) (*models.PaginatedResponse, error) {
	// Verificar permissão
	if err := s.pacienteValidator.ValidateAccess(ctx, 0, "read"); err != nil {
		return nil, err
	}
	
	// Validar parâmetros de paginação
	if filtro.Page < 1 {
		filtro.Page = 1
	}
	if filtro.PageSize < 1 || filtro.PageSize > 100 {
		filtro.PageSize = 10
	}
	
	// Buscar pacientes
	pacientes, total, err := s.pacienteRepo.Search(ctx, filtro)
	if err != nil {
		return nil, err
	}
	
	// Converter modelos para respostas
	var responses []models.PacienteResponse
	for i := range pacientes {
		responses = append(responses, *s.dtoConverter.ModelToResponse(&pacientes[i]))
	}
	
	// Criar resposta paginada
	response := &models.PaginatedResponse{
		Total:    total,
		Page:     filtro.Page,
		PageSize: filtro.PageSize,
		Data:     responses,
	}
	
	return response, nil
}