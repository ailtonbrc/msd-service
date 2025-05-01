package repository

import (
	"context"
	"errors"

	"clinica_server/internal/models"

	"gorm.io/gorm"
)

// PacienteRepository define a interface para operações de repositório de pacientes
type PacienteRepository interface {
	// Create cria um novo paciente no banco de dados
	Create(ctx context.Context, paciente *models.Paciente) error
	
	// GetByID busca um paciente pelo ID
	GetByID(ctx context.Context, id uint) (*models.Paciente, error)
	
	// GetAll retorna todos os pacientes com paginação e busca opcional
	GetAll(ctx context.Context, page, pageSize int, search string) ([]models.Paciente, int64, error)
	
	// Update atualiza um paciente existente
	Update(ctx context.Context, paciente *models.Paciente) error
	
	// Delete exclui um paciente pelo ID (soft delete)
	Delete(ctx context.Context, id uint) error
	
	// ExistsByCPF verifica se existe um paciente com o CPF informado
	ExistsByCPF(ctx context.Context, cpf string, excludeID uint) (bool, error)
	
	// Search busca pacientes com filtros avançados
	Search(ctx context.Context, filtro *models.PacienteFiltro) ([]models.Paciente, int64, error)
}

// GormPacienteRepository implementa PacienteRepository usando GORM
type GormPacienteRepository struct {
	db *gorm.DB
}

// NewPacienteRepository cria uma nova instância de PacienteRepository
func NewPacienteRepository(db *gorm.DB) PacienteRepository {
	return &GormPacienteRepository{db: db}
}

// Create cria um novo paciente no banco de dados
// Recebe o contexto e o modelo do paciente a ser criado
// Retorna erro em caso de falha
func (r *GormPacienteRepository) Create(ctx context.Context, paciente *models.Paciente) error {
	// Usa o contexto para garantir cancelamento adequado da operação
	return r.db.WithContext(ctx).Create(paciente).Error
}

// GetByID busca um paciente pelo ID
// Recebe o contexto e o ID do paciente
// Retorna o paciente encontrado ou nil se não existir, e erro em caso de falha
func (r *GormPacienteRepository) GetByID(ctx context.Context, id uint) (*models.Paciente, error) {
	var paciente models.Paciente
	
	// Busca o paciente pelo ID
	result := r.db.WithContext(ctx).First(&paciente, id)
	
	// Verifica se houve erro
	if result.Error != nil {
		// Se o erro for "registro não encontrado", retorna nil sem erro
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// Caso contrário, retorna o erro
		return nil, result.Error
	}
	
	return &paciente, nil
}

// GetAll retorna todos os pacientes com paginação e busca opcional
// Recebe o contexto, página, tamanho da página e termo de busca
// Retorna a lista de pacientes, o total de registros e erro em caso de falha
func (r *GormPacienteRepository) GetAll(ctx context.Context, page, pageSize int, search string) ([]models.Paciente, int64, error) {
	var pacientes []models.Paciente
	var total int64

	// Inicia a query
	query := r.db.WithContext(ctx).Model(&models.Paciente{})

	// Aplica filtro de busca se fornecido
	if search != "" {
		query = query.Where("nome LIKE ? OR cpf LIKE ? OR email LIKE ? OR telefone LIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Conta o total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calcula o offset para paginação
	offset := (page - 1) * pageSize
	
	// Busca os pacientes com paginação e ordenação
	if err := query.Offset(offset).Limit(pageSize).Order("nome ASC").Find(&pacientes).Error; err != nil {
		return nil, 0, err
	}

	return pacientes, total, nil
}

// Update atualiza um paciente existente
// Recebe o contexto e o modelo do paciente com as atualizações
// Retorna erro em caso de falha
func (r *GormPacienteRepository) Update(ctx context.Context, paciente *models.Paciente) error {
	return r.db.WithContext(ctx).Save(paciente).Error
}

// Delete exclui um paciente pelo ID (soft delete)
// Recebe o contexto e o ID do paciente a ser excluído
// Retorna erro em caso de falha
func (r *GormPacienteRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Paciente{}, id).Error
}

// ExistsByCPF verifica se existe um paciente com o CPF informado
// Recebe o contexto, o CPF a ser verificado e um ID opcional para exclusão da verificação
// Retorna true se existir, false caso contrário, e erro em caso de falha
func (r *GormPacienteRepository) ExistsByCPF(ctx context.Context, cpf string, excludeID uint) (bool, error) {
	var count int64
	
	// Inicia a query
	query := r.db.WithContext(ctx).Model(&models.Paciente{}).Where("cpf = ?", cpf)
	
	// Se excludeID for maior que 0, exclui esse paciente da verificação
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	
	// Conta os registros
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// Search busca pacientes com filtros avançados
// Recebe o contexto e os filtros de busca
// Retorna a lista de pacientes, o total de registros e erro em caso de falha
func (r *GormPacienteRepository) Search(ctx context.Context, filtro *models.PacienteFiltro) ([]models.Paciente, int64, error) {
	var pacientes []models.Paciente
	var total int64

	// Inicia a query
	query := r.db.WithContext(ctx).Model(&models.Paciente{})

	// Aplica os filtros
	if filtro.Nome != "" {
		query = query.Where("nome LIKE ?", "%"+filtro.Nome+"%")
	}
	
	if filtro.CPF != "" {
		query = query.Where("cpf LIKE ?", "%"+filtro.CPF+"%")
	}
	
	if filtro.Telefone != "" {
		query = query.Where("telefone LIKE ?", "%"+filtro.Telefone+"%")
	}
	
	if filtro.Diagnostico != "" {
		query = query.Where("diagnostico LIKE ?", "%"+filtro.Diagnostico+"%")
	}

	// Conta o total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calcula o offset para paginação
	offset := (filtro.Page - 1) * filtro.PageSize
	
	// Busca os pacientes com paginação e ordenação
	if err := query.Offset(offset).Limit(filtro.PageSize).Order("nome ASC").Find(&pacientes).Error; err != nil {
		return nil, 0, err
	}

	return pacientes, total, nil
}