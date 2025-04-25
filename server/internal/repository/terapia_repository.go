package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msd-service/server/internal/models"
)

// TerapiaRepository define a interface para operações de repositório de terapias
type TerapiaRepository interface {
	Create(ctx context.Context, terapia *models.Terapia) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Terapia, error)
	Update(ctx context.Context, terapia *models.Terapia) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Terapia, error)
	Count(ctx context.Context) (int64, error)
}

// GormTerapiaRepository implementa TerapiaRepository usando GORM
type GormTerapiaRepository struct {
	db *gorm.DB
}

// NewGormTerapiaRepository cria uma nova instância de GormTerapiaRepository
func NewGormTerapiaRepository(db *gorm.DB) *GormTerapiaRepository {
	return &GormTerapiaRepository{db: db}
}

// Create cria uma nova terapia no banco de dados
func (r *GormTerapiaRepository) Create(ctx context.Context, terapia *models.Terapia) error {
	return r.db.WithContext(ctx).Create(terapia).Error
}

// GetByID busca uma terapia pelo ID
func (r *GormTerapiaRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Terapia, error) {
	var terapia models.Terapia
	if err := r.db.WithContext(ctx).First(&terapia, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &terapia, nil
}

// Update atualiza uma terapia existente
func (r *GormTerapiaRepository) Update(ctx context.Context, terapia *models.Terapia) error {
	return r.db.WithContext(ctx).Save(terapia).Error
}

// Delete exclui uma terapia pelo ID (soft delete)
func (r *GormTerapiaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Terapia{}, "id = ?", id).Error
}

// List retorna uma lista paginada de terapias
func (r *GormTerapiaRepository) List(ctx context.Context, limit, offset int) ([]*models.Terapia, error) {
	var terapias []*models.Terapia
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&terapias).Error; err != nil {
		return nil, err
	}
	return terapias, nil
}

// Count retorna o número total de terapias
func (r *GormTerapiaRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Terapia{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
