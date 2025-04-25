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
	ErrSessaoNotFound = errors.New("sessão não encontrada")
)

// SessaoService encapsula a lógica de negócio relacionada a sessões
type SessaoService struct {
	repo repository.SessaoRepository
}

// NewSessaoService cria uma nova instância de SessaoService
func NewSessaoService(repo repository.SessaoRepository) *SessaoService {
	return &SessaoService{repo: repo}
}

// CreateSessao cria uma nova sessão
func (s *SessaoService) CreateSessao(ctx context.Context, sessao *models.Sessao) (*models.Sessao, error) {
	if err := s.repo.Create(ctx, sessao); err != nil {
		return nil, err
	}
	return sessao, nil
}

// GetSessao busca uma sessão pelo ID
func (s *SessaoService) GetSessao(ctx context.Context, id uuid.UUID) (*models.Sessao, error) {
	sessao, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sessao == nil {
		return nil, ErrSessaoNotFound
	}
	return sessao, nil
}

// UpdateSessao atualiza uma sessão existente
func (s *SessaoService) UpdateSessao(ctx context.Context, sessao *models.Sessao) (*models.Sessao, error) {
	existing, err := s.repo.GetByID(ctx, sessao.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrSessaoNotFound
	}

	if err := s.repo.Update(ctx, sessao); err != nil {
		return nil, err
	}
	return sessao, nil
}

// DeleteSessao exclui uma sessão pelo ID
func (s *SessaoService) DeleteSessao(ctx context.Context, id uuid.UUID) error {
	sessao, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if sessao == nil {
		return ErrSessaoNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ListSessoes retorna uma lista paginada de sessões
func (s *SessaoService) ListSessoes(ctx context.Context, page, pageSize int) ([]*models.Sessao, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	sessoes, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return sessoes, total, nil
}

// ListSessoesByPaciente retorna uma lista paginada de sessões de um paciente específico
func (s *SessaoService) ListSessoesByPaciente(ctx context.Context, pacienteID uuid.UUID, page, pageSize int) ([]*models.Sessao, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	sessoes, err := s.repo.ListByPaciente(ctx, pacienteID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByPaciente(ctx, pacienteID)
	if err != nil {
		return nil, 0, err
	}

	return sessoes, total, nil
}
