// internal/repository/paciente_repository.go
package repository

import (
	"context"
	"errors"
	"fmt"

	"clinica_server/internal/models"

	"gorm.io/gorm"
)

// Erros comuns do repositório
var (
	ErrPacienteNotFound = errors.New("paciente não encontrado")
	ErrInvalidID        = errors.New("ID inválido")
	ErrDuplicateCPF     = errors.New("CPF já cadastrado")
	ErrDatabase         = errors.New("erro no banco de dados")
)

// PacienteRepository define a interface para operações de repositório de pacientes
type PacienteRepository interface {
	// Operações básicas CRUD
	Create(ctx context.Context, paciente *models.Paciente) error
	GetByID(ctx context.Context, id uint) (*models.Paciente, error)
	Update(ctx context.Context, paciente *models.Paciente) error
	Delete(ctx context.Context, id uint, userID uint) error
	
	// Operações de listagem e busca
	List(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.Paciente, int64, error)
	GetByCPF(ctx context.Context, cpf string) (*models.Paciente, error)
	Search(ctx context.Context, query string, page, limit int) ([]models.Paciente, int64, error)
	
	// Verificações
	ExistsByCPF(ctx context.Context, cpf string, excludeID uint) (bool, error)
}

// GormPacienteRepository implementa PacienteRepository usando GORM
type GormPacienteRepository struct {
	db *gorm.DB
}

// NewPacienteRepository cria uma nova instância de PacienteRepository
func NewPacienteRepository(db *gorm.DB) PacienteRepository {
	return &GormPacienteRepository{
		db: db,
	}
}

// Create cria um novo paciente no banco de dados
func (r *GormPacienteRepository) Create(ctx context.Context, paciente *models.Paciente) error {
	// Verificar se já existe um paciente com o mesmo CPF
	if paciente.CPF != "" {
		exists, err := r.ExistsByCPF(ctx, paciente.CPF, 0)
		if err != nil {
			return fmt.Errorf("erro ao verificar CPF existente: %w", err)
		}
		if exists {
			return ErrDuplicateCPF
		}
	}

	// Usar o contexto para cancelamento
	tx := r.db.WithContext(ctx)
	
	// Criar o paciente
	if err := tx.Create(paciente).Error; err != nil {
		return fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	return nil
}

// GetByID busca um paciente pelo ID
func (r *GormPacienteRepository) GetByID(ctx context.Context, id uint) (*models.Paciente, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}
	
	var paciente models.Paciente
	tx := r.db.WithContext(ctx)
	
	if err := tx.First(&paciente, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPacienteNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	return &paciente, nil
}

// Update atualiza um paciente existente
func (r *GormPacienteRepository) Update(ctx context.Context, paciente *models.Paciente) error {
	if paciente.ID == 0 {
		return ErrInvalidID
	}
	
	// Verificar se o paciente existe
	_, err := r.GetByID(ctx, paciente.ID)
	if err != nil {
		return err
	}
	
	// Verificar se o CPF já está em uso por outro paciente
	if paciente.CPF != "" {
		exists, err := r.ExistsByCPF(ctx, paciente.CPF, paciente.ID)
		if err != nil {
			return fmt.Errorf("erro ao verificar CPF existente: %w", err)
		}
		if exists {
			return ErrDuplicateCPF
		}
	}
	
	tx := r.db.WithContext(ctx)
	
	// Atualizar apenas os campos não-zero
	if err := tx.Model(paciente).Updates(paciente).Error; err != nil {
		return fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	return nil
}

// Delete remove um paciente (soft delete devido ao DeletedAt no modelo)
func (r *GormPacienteRepository) Delete(ctx context.Context, id uint, userID uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	
	// Verificar se o paciente existe
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	tx := r.db.WithContext(ctx)
	
	// Soft delete com informação de quem excluiu
	if err := tx.Model(&models.Paciente{}).Where("id = ?", id).
		Updates(map[string]interface{}{"atualizado_por": userID}).
		Delete(&models.Paciente{}).Error; err != nil {
		return fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	return nil
}

// List retorna uma lista paginada de pacientes com filtros opcionais
func (r *GormPacienteRepository) List(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.Paciente, int64, error) {
	var pacientes []models.Paciente
	var total int64
	
	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	offset := (page - 1) * limit
	
	// Iniciar a query
	tx := r.db.WithContext(ctx).Model(&models.Paciente{})
	
	// Aplicar filtros
	for key, value := range filters {
		switch key {
		case "nome":
			tx = tx.Where("nome LIKE ?", "%"+value.(string)+"%")
		case "cpf":
			tx = tx.Where("cpf = ?", value)
		case "diagnostico":
			tx = tx.Where("diagnostico LIKE ?", "%"+value.(string)+"%")
		case "genero":
			tx = tx.Where("genero = ?", value)
		case "cidade":
			tx = tx.Where("cidade = ?", value)
		case "estado":
			tx = tx.Where("estado = ?", value)
		case "criado_por":
			tx = tx.Where("criado_por = ?", value)
		}
	}
	
	// Contar total de registros
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	// Executar a query com paginação
	if err := tx.Limit(limit).Offset(offset).Order("nome ASC").Find(&pacientes).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	return pacientes, total, nil
}

// GetByCPF busca um paciente pelo CPF
func (r *GormPacienteRepository) GetByCPF(ctx context.Context, cpf string) (*models.Paciente, error) {
	if cpf == "" {
		return nil, errors.New("CPF não pode ser vazio")
	}
	
	var paciente models.Paciente
	tx := r.db.WithContext(ctx)
	
	if err := tx.Where("cpf = ?", cpf).First(&paciente).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPacienteNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	return &paciente, nil
}

// Search busca pacientes por um termo de busca
func (r *GormPacienteRepository) Search(ctx context.Context, query string, page, limit int) ([]models.Paciente, int64, error) {
	var pacientes []models.Paciente
	var total int64
	
	// Validar parâmetros de paginação
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	
	offset := (page - 1) * limit
	
	// Construir a query de busca
	searchQuery := "%" + query + "%"
	tx := r.db.WithContext(ctx).Model(&models.Paciente{}).
		Where("nome LIKE ? OR cpf LIKE ? OR email LIKE ? OR telefone LIKE ? OR diagnostico LIKE ?",
			searchQuery, searchQuery, searchQuery, searchQuery, searchQuery)
	
	// Contar total de registros
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	// Executar a query com paginação
	if err := tx.Limit(limit).Offset(offset).Order("nome ASC").Find(&pacientes).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	return pacientes, total, nil
}

// ExistsByCPF verifica se existe um paciente com o CPF informado (excluindo o ID opcional)
func (r *GormPacienteRepository) ExistsByCPF(ctx context.Context, cpf string, excludeID uint) (bool, error) {
	if cpf == "" {
		return false, nil
	}
	
	var count int64
	tx := r.db.WithContext(ctx).Model(&models.Paciente{}).Where("cpf = ?", cpf)
	
	// Se excludeID for fornecido, exclui esse ID da busca
	if excludeID > 0 {
		tx = tx.Where("id != ?", excludeID)
	}
	
	if err := tx.Count(&count).Error; err != nil {
		return false, fmt.Errorf("%w: %v", ErrDatabase, err)
	}
	
	return count > 0, nil
}