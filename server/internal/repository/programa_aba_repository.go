package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msd-service/server/internal/models"
)

// ProgramaABARepository define a interface para operações de repositório de programas ABA
type ProgramaABARepository interface {
	Create(ctx context.Context, programa *models.ProgramaABA) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.ProgramaABA, error)
	Update(ctx context.Context, programa *models.ProgramaABA) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.ProgramaABA, error)
	ListByPaciente(ctx context.Context, pacienteID uuid.UUID, limit, offset int) ([]*models.ProgramaABA, error)
	Count(ctx context.Context) (int64, error)
	CountByPaciente(ctx context.Context, pacienteID uuid.UUID) (int64, error)
}

// GormProgramaABARepository implementa ProgramaABARepository usando GORM
type GormProgramaABARepository struct {
	db *gorm.DB
}

// NewGormProgramaABARepository cria uma nova instância de GormProgramaABARepository
func NewGormProgramaABARepository(db *gorm.DB) *GormProgramaABARepository {
	return &GormProgramaABARepository{db: db}
}

// Create cria um novo programa ABA no banco de dados
func (r *GormProgramaABARepository) Create(ctx context.Context, programa *models.ProgramaABA) error {
	return r.db.WithContext(ctx).Create(programa).Error
}

// GetByID busca um programa ABA pelo ID
func (r *GormProgramaABARepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ProgramaABA, error) {
	var programa models.ProgramaABA
	if err := r.db.WithContext(ctx).First(&programa, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &programa, nil
}

// Update atualiza um programa ABA existente
func (r *GormProgramaABARepository) Update(ctx context.Context, programa *models.ProgramaABA) error {
	return r.db.WithContext(ctx).Save(programa).Error
}

// Delete exclui um programa ABA pelo ID (soft delete)
func (r *GormProgramaABARepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ProgramaABA{}, "id = ?", id).Error
}

// List retorna uma lista paginada de programas ABA
func (r *GormProgramaABARepository) List(ctx context.Context, limit, offset int) ([]*models.ProgramaABA, error) {
	var programas []*models.ProgramaABA
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&programas).Error; err != nil {
		return nil, err
	}
	return programas, nil
}

// ListByPaciente retorna uma lista paginada de programas ABA de um paciente específico
func (r *GormProgramaABARepository) ListByPaciente(ctx context.Context, pacienteID uuid.UUID, limit, offset int) ([]*models.ProgramaABA, error) {
	var programas []*models.ProgramaABA
	if err := r.db.WithContext(ctx).Where("paciente_id = ?", pacienteID).Limit(limit).Offset(offset).Find(&programas).Error; err != nil {
		return nil, err
	}
	return programas, nil
}

// Count retorna o número total de programas ABA
func (r *GormProgramaABARepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.ProgramaABA{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByPaciente retorna o número total de programas ABA de um paciente específico
func (r *GormProgramaABARepository) CountByPaciente(ctx context.Context, pacienteID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.ProgramaABA{}).Where("paciente_id = ?", pacienteID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
