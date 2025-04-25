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
	ErrProgramaNotFound = errors.New("programa ABA não encontrado")
)

// ProgramaABAService encapsula a lógica de negócio relacionada a programas ABA
type ProgramaABAService struct {
	repo repository.ProgramaABARepository
}

// NewProgramaABAService cria uma nova instância de ProgramaABAService
func NewProgramaABAService(repo repository.ProgramaABARepository) *ProgramaABAService {
	return &ProgramaABAService{repo: repo}
}

// CreatePrograma cria um novo programa ABA
func (s *ProgramaABAService) CreatePrograma(ctx context.Context, programa *models.ProgramaABA) (*models.ProgramaABA, error) {
	if err := s.repo.Create(ctx, programa); err != nil {
		return nil, err
	}
	return programa, nil
}

// GetPrograma busca um programa ABA pelo ID
func (s *ProgramaABAService) GetPrograma(ctx context.Context, id uuid.UUID) (*models.ProgramaABA, error) {
	programa, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if programa == nil {
		return nil, ErrProgramaNotFound
	}
	return programa, nil
}

// UpdatePrograma atualiza um programa ABA existente
func (s *ProgramaABAService) UpdatePrograma(ctx context.Context, programa *models.ProgramaABA) (*models.ProgramaABA, error) {
	existing, err := s.repo.GetByID(ctx, programa.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrProgramaNotFound
	}

	if err := s.repo.Update(ctx, programa); err != nil {
		return nil, err
	}
	return programa, nil
}

// DeletePrograma exclui um programa ABA pelo ID
func (s *ProgramaABAService) DeletePrograma(ctx context.Context, id uuid.UUID) error {
	programa, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if programa == nil {
		return ErrProgramaNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ListProgramas retorna uma lista paginada de programas ABA
func (s *ProgramaABAService) ListProgramas(ctx context.Context, page, pageSize int) ([]*models.ProgramaABA, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	programas, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return programas, total, nil
}

// ListProgramasByPaciente retorna uma lista paginada de programas ABA de um paciente específico
func (s *ProgramaABAService) ListProgramasByPaciente(ctx context.Context, pacienteID uuid.UUID, page, pageSize int) ([]*models.ProgramaABA, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	programas, err := s.repo.ListByPaciente(ctx, pacienteID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByPaciente(ctx, pacienteID)
	if err != nil {
		return nil, 0, err
	}

	return programas, total, nil
}
