package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msd-service/server/internal/models"
)

// ObjetivoTerapeuticoRepository define a interface para operações de repositório de objetivos terapêuticos
type ObjetivoTerapeuticoRepository interface {
	Create(ctx context.Context, objetivo *models.ObjetivoTerapeutico) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ObjetivoTerapeutico, error)
	Update(ctx context.Context, objetivo *models.ObjetivoTerapeutico) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.ObjetivoTerapeutico, error)
	ListByPaciente(ctx context.Context, pacienteID uuid.UUID, limit, offset int) ([]*models.ObjetivoTerapeutico, error)
	Count(ctx context.Context) (int64, error)
	CountByPaciente(ctx context.Context, pacienteID uuid.UUID) (int64, error)
}

// GormObjetivoTerapeuticoRepository implementa ObjetivoTerapeuticoRepository usando GORM
type GormObjetivoTerapeuticoRepository struct {
	db *gorm.DB
}

// NewGormObjetivoTerapeuticoRepository cria uma nova instância de GormObjetivoTerapeuticoRepository
func NewGormObjetivoTerapeuticoRepository(db *gorm.DB) *GormObjetivoTerapeuticoRepository {
	return &GormObjetivoTerapeuticoRepository{db: db}
}

// Create cria um novo objetivo terapêutico no banco de dados
func (r *GormObjetivoTerapeuticoRepository) Create(ctx context.Context, objetivo *models.ObjetivoTerapeutico) error {
	return r.db.WithContext(ctx).Create(objetivo).Error
}

// GetByID busca um objetivo terapêutico pelo ID
func (r *GormObjetivoTerapeuticoRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ObjetivoTerapeutico, error) {
	var objetivo models.ObjetivoTerapeutico
	if err := r.db.WithContext(ctx).First(&objetivo, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &objetivo, nil
}

// Update atualiza um objetivo terapêutico existente
func (r *GormObjetivoTerapeuticoRepository) Update(ctx context.Context, objetivo *models.ObjetivoTerapeutico) error {
	return r.db.WithContext(ctx).Save(objetivo).Error
}

// Delete exclui um objetivo terapêutico pelo ID (soft delete)
func (r *GormObjetivoTerapeuticoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ObjetivoTerapeutico{}, "id = ?", id).Error
}

// List retorna uma lista paginada de objetivos terapêuticos
func (r *GormObjetivoTerapeuticoRepository) List(ctx context.Context, limit, offset int) ([]*models.ObjetivoTerapeutico, error) {
	var objetivos []*models.ObjetivoTerapeutico
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&objetivos).Error; err != nil {
		return nil, err
	}
	return objetivos, nil
}

// ListByPaciente retorna uma lista paginada de objetivos terapêuticos de um paciente específico
func (r *GormObjetivoTerapeuticoRepository) ListByPaciente(ctx context.Context, pacienteID uuid.UUID, limit, offset int) ([]*models.ObjetivoTerapeutico, error) {
	var objetivos []*models.ObjetivoTerapeutico
	if err := r.db.WithContext(ctx).Where("paciente_id = ?", pacienteID).Limit(limit).Offset(offset).Find(&objetivos).Error; err != nil {
		return nil, err
	}
	return objetivos, nil
}

// Count retorna o número total de objetivos terapêuticos
func (r *GormObjetivoTerapeuticoRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.ObjetivoTerapeutico{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByPaciente retorna o número total de objetivos terapêuticos de um paciente específico
func (r *GormObjetivoTerapeuticoRepository) CountByPaciente(ctx context.Context, pacienteID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.ObjetivoTerapeutico{}).Where("paciente_id = ?", pacienteID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
