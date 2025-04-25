package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"msd-service/server/internal/models"
	"msd-service/server/internal/repository"
)

// Erros comuns do serviço
var (
	ErrPacienteNotFound = errors.New("paciente não encontrado")
	ErrInvalidInput     = errors.New("dados de entrada inválidos")
)

// PacienteService encapsula a lógica de negócio relacionada a pacientes
type PacienteService struct {
	repo repository.PacienteRepository
}

// NewPacienteService cria uma nova instância de PacienteService
func NewPacienteService(repo repository.PacienteRepository) *PacienteService {
	return &PacienteService{repo: repo}
}

// CreatePaciente cria um novo paciente
func (s *PacienteService) CreatePaciente(ctx context.Context, req *models.CreatePacienteRequest) (*models.PacienteResponse, error) {
	if req == nil {
		return nil, ErrInvalidInput
	}

	paciente := req.ToPaciente()
	if err := s.repo.Create(ctx, paciente); err != nil {
		return nil, err
	}

	return paciente.ToResponse(), nil
}

// GetPaciente busca um paciente pelo ID
func (s *PacienteService) GetPaciente(ctx context.Context, id uuid.UUID) (*models.PacienteResponse, error) {
	paciente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if paciente == nil {
		return nil, ErrPacienteNotFound
	}

	return paciente.ToResponse(), nil
}

// UpdatePaciente atualiza um paciente existente
func (s *PacienteService) UpdatePaciente(ctx context.Context, id uuid.UUID, req *models.UpdatePacienteRequest) (*models.PacienteResponse, error) {
	if req == nil {
		return nil, ErrInvalidInput
	}

	paciente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if paciente == nil {
		return nil, ErrPacienteNotFound
	}

	paciente.ApplyUpdates(req)
	if err := s.repo.Update(ctx, paciente); err != nil {
		return nil, err
	}

	return paciente.ToResponse(), nil
}

// DeletePaciente exclui um paciente pelo ID
func (s *PacienteService) DeletePaciente(ctx context.Context, id uuid.UUID) error {
	paciente, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if paciente == nil {
		return ErrPacienteNotFound
	}

	return s.repo.Delete(ctx, id)
}

// ListPacientes retorna uma lista paginada de pacientes
func (s *PacienteService) ListPacientes(ctx context.Context, page, pageSize int) ([]*models.PacienteResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	pacientes, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*models.PacienteResponse, len(pacientes))
	for i, paciente := range pacientes {
		responses[i] = paciente.ToResponse()
	}

	return responses, total, nil
}
