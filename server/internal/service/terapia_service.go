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
	ErrTerapiaNotFound = errors.New("terapia não encontrada")
)

// TerapiaService encapsula a lógica de negócio relacionada a terapias
type TerapiaService struct {
	repo repository.TerapiaRepository
}

// NewTerapiaService cria uma nova instância de TerapiaService
func NewTerapiaService(repo repository.TerapiaRepository) *TerapiaService {
	return &TerapiaService{repo: repo}
}

// CreateTerapia cria uma nova terapia
func (s *TerapiaService) CreateTerapia(ctx context.Context, terapia *models.Terapia) (*models.Terapia, error) {
	if err := s.repo.Create(ctx, terapia); err != nil {
		return nil, err
	}
	return terapia, nil
}

// GetTerapia busca uma terapia pelo ID
func (s *TerapiaService) GetTerapia(ctx context.Context, id uuid.UUID) (*models.Terapia, error) {
	terapia, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if terapia == nil {
		return nil, ErrTerapiaNotFound
	}
	return terapia, nil
}

// UpdateTerapia atualiza uma terapia existente
func (s *TerapiaService) UpdateTerapia(ctx context.Context, terapia *models.Terapia) (*models.Terapia, error) {
	existing, err := s.repo.GetByID(ctx, terapia.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrTerapiaNotFound
	}

	if err := s.repo.Update(ctx, terapia); err != nil {
		return nil, err
	}
	return terapia, nil
}

// DeleteTerapia exclui uma terapia pelo ID
func (s *TerapiaService) DeleteTerapia(ctx context.Context, id uuid.UUID) error {
	terapia, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if terapia == nil {
		return ErrTerapiaNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ListTerapias retorna uma lista paginada de terapias
func (s *TerapiaService) ListTerapias(ctx context.Context, page, pageSize int) ([]*models.Terapia, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	terapias, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return terapias, total, nil
}
