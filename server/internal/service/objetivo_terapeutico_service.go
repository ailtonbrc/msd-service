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
	ErrObjetivoNotFound = errors.New("objetivo terapêutico não encontrado")
)

// ObjetivoTerapeuticoService encapsula a lógica de negócio relacionada a objetivos terapêuticos
type ObjetivoTerapeuticoService struct {
	repo repository.ObjetivoTerapeuticoRepository
}

// NewObjetivoTerapeuticoService cria uma nova instância de ObjetivoTerapeuticoService
func NewObjetivoTerapeuticoService(repo repository.ObjetivoTerapeuticoRepository) *ObjetivoTerapeuticoService {
	return &ObjetivoTerapeuticoService{repo: repo}
}

// CreateObjetivo cria um novo objetivo terapêutico
func (s *ObjetivoTerapeuticoService) CreateObjetivo(ctx context.Context, objetivo *models.ObjetivoTerapeutico) (*models.ObjetivoTerapeutico, error) {
	if err := s.repo.Create(ctx, objetivo); err != nil {
		return nil, err
	}
	return objetivo, nil
}

// GetObjetivo busca um objetivo terapêutico pelo ID
func (s *ObjetivoTerapeuticoService) GetObjetivo(ctx context.Context, id uuid.UUID) (*models.ObjetivoTerapeutico, error) {
	objetivo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if objetivo == nil {
		return nil, ErrObjetivoNotFound
	}
	return objetivo, nil
}

// UpdateObjetivo atualiza um objetivo terapêutico existente
func (s *ObjetivoTerapeuticoService) UpdateObjetivo(ctx context.Context, objetivo *models.ObjetivoTerapeutico) (*models.ObjetivoTerapeutico, error) {
	existing, err := s.repo.GetByID(ctx, objetivo.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrObjetivoNotFound
	}

	if err := s.repo.Update(ctx, objetivo); err != nil {
		return nil, err
	}
	return objetivo, nil
}

// DeleteObjetivo exclui um objetivo terapêutico pelo ID
func (s *ObjetivoTerapeuticoService) DeleteObjetivo(ctx context.Context, id uuid.UUID) error {
	objetivo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if objetivo == nil {
		return ErrObjetivoNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ListObjetivos retorna uma lista paginada de objetivos terapêuticos
func (s *ObjetivoTerapeuticoService) ListObjetivos(ctx context.Context, page, pageSize int) ([]*models.ObjetivoTerapeutico, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	objetivos, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return objetivos, total, nil
}

// ListObjetivosByPaciente retorna uma lista paginada de objetivos terapêuticos de um paciente específico
func (s *ObjetivoTerapeuticoService) ListObjetivosByPaciente(ctx context.Context, pacienteID uuid.UUID, page, pageSize int) ([]*models.ObjetivoTerapeutico, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	objetivos, err := s.repo.ListByPaciente(ctx, pacienteID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByPaciente(ctx, pacienteID)
	if err != nil {
		return nil, 0, err
	}

	return objetivos, total, nil
}
