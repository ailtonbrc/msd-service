package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msd-service/server/internal/models"
)

// ComportamentoAlvoRepository define a interface para operações de repositório de comportamentos alvo
type ComportamentoAlvoRepository interface {
	Create(ctx context.Context, comportamento *models.ComportamentoAlvo) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ComportamentoAlvo  error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ComportamentoAlvo, error)
	Update(ctx context.Context, comportamento *models.ComportamentoAlvo) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.ComportamentoAlvo, error)
	ListByPaciente(ctx context.Context, pacienteID uuid.UUID, limit, offset int) ([]*models.ComportamentoAlvo, error)
	Count(ctx context.Context) (int64, error)
	CountByPaciente(ctx context.Context, pacienteID uuid.UUID) (int64, error)
}

// GormComportamentoAlvoRepository implementa ComportamentoAlvoRepository usando GORM
type GormComportamentoAlvoRepository struct {
	db *gorm.DB
}

// NewGormComportamentoAlvoRepository cria uma nova instância de GormComportamentoAlvoRepository
func NewGormComportamentoAlvoRepository(db *gorm.DB) *GormComportamentoAlvoRepository {
	return &GormComportamentoAlvoRepository{db: db}
}

// Create cria um novo comportamento alvo no banco de dados
func (r *GormComportamentoAlvoRepository) Create(ctx context.Context, comportamento *models.ComportamentoAlvo) error {
	return r.db.WithContext(ctx).Create(comportamento).Error
}

// GetByID busca um comportamento alvo pelo ID
func (r *GormComportamentoAlvoRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ComportamentoAlvo, error) {
	var comportamento models.ComportamentoAlvo
	if err := r.db.WithContext(ctx).First(&comportamento, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comportamento, nil
}

// Update atualiza um comportamento alvo existente
func (r *GormComportamentoAlvoRepository) Update(ctx context.Context, comportamento *models.ComportamentoAlvo) error {
	return r.db.WithContext(ctx).Save(comportamento).Error
}

// Delete exclui um comportamento alvo pelo ID (soft delete)
func (r *GormComportamentoAlvoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ComportamentoAlvo{}, "id = ?", id).Error
}

// List retorna uma lista paginada de comportamentos alvo
func (r *GormComportamentoAlvoRepository) List(ctx context.Context, limit, offset int) ([]*models.ComportamentoAlvo, error) {
	var comportamentos []*models.ComportamentoAlvo
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&comportamentos).Error; err != nil {
		return nil, err
	}
	return comportamentos, nil
}

// ListByPaciente retorna uma lista paginada de comportamentos alvo de um paciente específico
func (r *GormComportamentoAlvoRepository) ListByPaciente(ctx context.Context, pacienteID uuid.UUID, limit, offset int) ([]*models.ComportamentoAlvo, error) {
	var comportamentos []*models.ComportamentoAlvo
	if err := r.db.WithContext(ctx).Where("paciente_id = ?", pacienteID).Limit(limit).Offset(offset).Find(&comportamentos).Error; err != nil {
		return nil, err
	}
	return comportamentos, nil
}

// Count retorna o número total de comportamentos alvo
func (r *GormComportamentoAlvoRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.ComportamentoAlvo{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByPaciente retorna o número total de comportamentos alvo de um paciente específico
func (r *GormComportamentoAlvoRepository) CountByPaciente(ctx context.Context, pacienteID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.ComportamentoAlvo{}).Where("paciente_id = ?", pacienteID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
