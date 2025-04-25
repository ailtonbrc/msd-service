package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"msd-service/server/internal/models"
)

// PacienteRepository define a interface para operações de repositório de pacientes
type PacienteRepository interface {
	Create(ctx context.Context, paciente *models.Paciente) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Paciente, error)
	Update(ctx context.Context, paciente *models.Paciente) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Paciente, error)
	Count(ctx context.Context) (int64, error)
}

// GormPacienteRepository implementa PacienteRepository usando GORM
type GormPacienteRepository struct {
	db *gorm.DB
}

// NewGormPacienteRepository cria uma nova instância de GormPacienteRepository
func NewGormPacienteRepository(db *gorm.DB) *GormPacienteRepository {
	return &GormPacienteRepository{db: db}
}

// Create cria um novo paciente no banco de dados
func (r *GormPacienteRepository) Create(ctx context.Context, paciente *models.Paciente) error {
	return r.db.WithContext(ctx).Create(paciente).Error
}

// GetByID busca um paciente pelo ID
func (r *GormPacienteRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Paciente, error) {
	var paciente models.Paciente
	if err := r.db.WithContext(ctx).First(&paciente, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Retorna nil, nil quando não encontra o registro
		}
		return nil, err
	}
	return &paciente, nil
}

// Update atualiza um paciente existente
func (r *GormPacienteRepository) Update(ctx context.Context, paciente *models.Paciente) error {
	return r.db.WithContext(ctx).Save(paciente).Error
}

// Delete exclui um paciente pelo ID (soft delete)
func (r *GormPacienteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Paciente{}, "id = ?", id).Error
}

// List retorna uma lista paginada de pacientes
func (r *GormPacienteRepository) List(ctx context.Context, limit, offset int) ([]*models.Paciente, error) {
	var pacientes []*models.Paciente
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&pacientes).Error; err != nil {
		return nil, err
	}
	return pacientes, nil
}

// Count retorna o número total de pacientes
func (r *GormPacienteRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Paciente{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
