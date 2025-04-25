package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msd-service/server/internal/models"
)

// SessaoRepository define a interface para operações de repositório de sessões
type SessaoRepository interface {
	Create(ctx context.Context, sessao *models.Sessao) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Sessao, error)
	Update(ctx context.Context, sessao *models.Sessao) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Sessao, error)
	ListByPaciente(ctx context.Context, pacienteID uuid.UUID, limit, offset int) ([]*models.Sessao, error)
	Count(ctx context.Context) (int64, error)
	CountByPaciente(ctx context.Context, pacienteID uuid.UUID) (int64, error)
}

// GormSessaoRepository implementa SessaoRepository usando GORM
type GormSessaoRepository struct {
	db *gorm.DB
}

// NewGormSessaoRepository cria uma nova instância de GormSessaoRepository
func NewGormSessaoRepository(db *gorm.DB) *GormSessaoRepository {
	return &GormSessaoRepository{db: db}
}

// Create cria uma nova sessão no banco de dados
func (r *GormSessaoRepository) Create(ctx context.Context, sessao *models.Sessao) error {
	return r.db.WithContext(ctx).Create(sessao).Error
}

// GetByID busca uma sessão pelo ID
func (r *GormSessaoRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Sessao, error) {
	var sessao models.Sessao
	if err := r.db.WithContext(ctx).First(&sessao, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sessao, nil
}

// Update atualiza uma sessão existente
func (r *GormSessaoRepository) Update(ctx context.Context, sessao *models.Sessao) error {
	return r.db.WithContext(ctx).Save(sessao).Error
}

// Delete exclui uma sessão pelo ID (soft delete)
func (r *GormSessaoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Sessao{}, "id = ?", id).Error
}

// List retorna uma lista paginada de sessões
func (r *GormSessaoRepository) List(ctx context.Context, limit, offset int) ([]*models.Sessao, error) {
	var sessoes []*models.Sessao
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&sessoes).Error; err != nil {
		return nil, err
	}
	return sessoes, nil
}

// ListByPaciente retorna uma lista paginada de sessões de um paciente específico
func (r *GormSessaoRepository) ListByPaciente(ctx context.Context, pacienteID uuid.UUID, limit, offset int) ([]*models.Sessao, error) {
	var sessoes []*models.Sessao
	if err := r.db.WithContext(ctx).Where("paciente_id = ?", pacienteID).Limit(limit).Offset(offset).Find(&sessoes).Error; err != nil {
		return nil, err
	}
	return sessoes, nil
}

// Count retorna o número total de sessões
func (r *GormSessaoRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Sessao{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByPaciente retorna o número total de sessões de um paciente específico
func (r *GormSessaoRepository) CountByPaciente(ctx context.Context, pacienteID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Sessao{}).Where("paciente_id = ?", pacienteID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
