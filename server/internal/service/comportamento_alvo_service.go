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
	ErrComportamentoNotFound = errors.New("comportamento alvo não encontrado")
)

// ComportamentoAlvoService encapsula a lógica de negócio relacionada a comportamentos alvo
type ComportamentoAlvoService struct {
	repo repository.ComportamentoAlvoRepository
}

// NewComportamentoAlvoService cria uma nova instância de ComportamentoAlvoService
func NewComportamentoAlvoService(repo repository.ComportamentoAlvoRepository) *ComportamentoAlvoService {
	return &ComportamentoAlvoService{repo: repo}
}

// CreateComportamento cria um novo comportamento alvo
func (s *ComportamentoAlvoService) CreateComportamento(ctx context.Context, comportamento *models.ComportamentoAlvo) (*models.ComportamentoAlvo, error) {
	if err := s.repo.Create(ctx, comportamento); err != nil {
		return nil, err
	}
	return comportamento, nil
}

// GetComportamento busca um comportamento alvo pelo ID
func (s *ComportamentoAlvoService) GetComportamento(ctx context.Context, id uuid.UUID) (*models.ComportamentoAlvo, error) {
	comportamento, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if comportamento == nil {
		return nil, ErrComportamentoNotFound
	}
	return comportamento, nil
}

// UpdateComportamento atualiza um comportamento alvo existente
func (s *ComportamentoAlvoService) UpdateComportamento(ctx context.Context, comportamento *models.ComportamentoAlvo) (*models.ComportamentoAlvo, error) {
	existing, err := s.repo.GetByID(ctx, comportamento.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrComportamentoNotFound
	}

	if err := s.repo.Update(ctx, comportamento); err != nil {
		return nil, err
	}
	return comportamento, nil
}

// DeleteComportamento exclui um comportamento alvo pelo ID
func (s *ComportamentoAlvoService) DeleteComportamento(ctx context.Context, id uuid.UUID) error {
	comportamento, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if comportamento == nil {
		return ErrComportamentoNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ListComportamentos retorna uma lista paginada de comportamentos alvo
func (s *ComportamentoAlvoService) ListComportamentos(ctx context.Context, page, pageSize int) ([]*models.ComportamentoAlvo, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	comportamentos, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return comportamentos, total, nil
}

// ListComportamentosByPaciente retorna uma lista paginada de comportamentos alvo de um paciente específico
func (s *ComportamentoAlvoService) ListComportamentosByPaciente(ctx context.Context, pacienteID uuid.UUID, page, pageSize int) ([]*models.ComportamentoAlvo, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	comportamentos, err := s.repo.ListByPaciente(ctx, pacienteID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByPaciente(ctx, pacienteID)
	if err != nil {
		return nil, 0, err
	}

	return comportamentos, total, nil
}
