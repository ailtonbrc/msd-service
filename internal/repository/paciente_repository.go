package repository

import (
	"clinica_server/internal/models"
	"clinica_server/internal/utils"
	"errors"

	"gorm.io/gorm"
)

// PacienteRepository define as operações de acesso a dados para usuários
type PacienteRepository interface {
	Repository
	BuscaTodos(pagination *utils.Pagination) ([]models.Paciente, error)
	BuscaPorID(id uint) (*models.Paciente, error)
	BuscaPorNome(nome string) (*models.Paciente, error)
	BuscaPorEmail(email string) (*models.Paciente, error)
	VerificaOutroPacientePorEmail(email string, id uint) (bool, error)
	Create(Paciente *models.Paciente) error
	Update(Paciente *models.Paciente) error
	Delete(id uint) error
}

// GormPacienteRepository implementa PacienteRepository usando GORM
type GormPacienteRepository struct {
	*BaseRepository
}

// NewPacienteRepository cria um novo repository de usuários
func NewPacienteRepository(db *gorm.DB) PacienteRepository {
	return &GormPacienteRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// BuscaTodos retorna todos os usuários com paginação
func (r *GormPacienteRepository) BuscaTodos(pagination *utils.Pagination) ([]models.Paciente, error) {
	var Pacientes []models.Paciente

	query := r.GetDB().Model(&models.Paciente{})
	query, err := utils.Paginate(&models.Paciente{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&Pacientes).Error; err != nil {
		return nil, err
	}

	return Pacientes, nil
}

// BuscaPorID busca um usuário pelo ID
func (r *GormPacienteRepository) BuscaPorID(id uint) (*models.Paciente, error) {
	var Paciente models.Paciente
	if err := r.GetDB().First(&Paciente, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &Paciente, nil
}

// BuscaPorNome busca um usuário pelo nome de usuário
func (r *GormPacienteRepository) BuscaPorNome(nome string) (*models.Paciente, error) {
	var Paciente models.Paciente
	if err := r.GetDB().Where("LOWER(nome) = LOWER(?)", nome).First(&Paciente).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &Paciente, nil
}

// BuscaPorEmail busca um usuário pelo email
func (r *GormPacienteRepository) BuscaPorEmail(email string) (*models.Paciente, error) {
	var Paciente models.Paciente
	if err := r.GetDB().Where("LOWER(email) = LOWER(?)", email).First(&Paciente).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &Paciente, nil
}

// ExistsByEmailExcept verifica se existe um usuário com o email especificado, exceto o usuário com o ID especificado
func (r *GormPacienteRepository) VerificaOutroPacientePorEmail(email string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Paciente{}).Where("LOWER(email) = LOWER(?) AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}

// Create cria um novo usuário
func (r *GormPacienteRepository) Create(Paciente *models.Paciente) error {
	return r.GetDB().Create(Paciente).Error
}

// Update atualiza um usuário existente
func (r *GormPacienteRepository) Update(Paciente *models.Paciente) error {
	return r.GetDB().Save(Paciente).Error
}

// Delete exclui um usuário (soft delete)
func (r *GormPacienteRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.Paciente{}, id).Error
}
