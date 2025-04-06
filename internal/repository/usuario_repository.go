package repository

import (
	"clinica_server/internal/models"
	"clinica_server/internal/utils"
	"errors"

	"gorm.io/gorm"
)

// UsuarioRepository define as operações de acesso a dados para usuários
type UsuarioRepository interface {
	Repository
	BuscaTodos(pagination *utils.Pagination) ([]models.Usuario, error)
	BuscaPorID(id uint) (*models.Usuario, error)
	BuscaPorNome(nome string) (*models.Usuario, error)
	BuscaPorEmail(email string) (*models.Usuario, error)
	VerificaOutroUsuarioPorEmail(email string, id uint) (bool, error)
	Create(Usuario *models.Usuario) error
	Update(Usuario *models.Usuario) error
	Delete(id uint) error
}

// GormUsuarioRepository implementa UsuarioRepository usando GORM
type GormUsuarioRepository struct {
	*BaseRepository
}

// NewUsuarioRepository cria um novo repository de usuários
func NewUsuarioRepository(db *gorm.DB) UsuarioRepository {
	return &GormUsuarioRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// BuscaTodos retorna todos os usuários com paginação
func (r *GormUsuarioRepository) BuscaTodos(pagination *utils.Pagination) ([]models.Usuario, error) {
	var Usuarios []models.Usuario

	query := r.GetDB().Model(&models.Usuario{})
	query, err := utils.Paginate(&models.Usuario{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&Usuarios).Error; err != nil {
		return nil, err
	}

	return Usuarios, nil
}

// BuscaPorID busca um usuário pelo ID
func (r *GormUsuarioRepository) BuscaPorID(id uint) (*models.Usuario, error) {
	var Usuario models.Usuario
	if err := r.GetDB().First(&Usuario, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &Usuario, nil
}

// BuscaPorNome busca um usuário pelo nome de usuário
func (r *GormUsuarioRepository) BuscaPorNome(nome string) (*models.Usuario, error) {
	var Usuario models.Usuario
	if err := r.GetDB().Where("LOWER(nome) = LOWER(?)", nome).First(&Usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &Usuario, nil
}

// BuscaPorEmail busca um usuário pelo email
func (r *GormUsuarioRepository) BuscaPorEmail(email string) (*models.Usuario, error) {
	var Usuario models.Usuario
	if err := r.GetDB().Where("LOWER(email) = LOWER(?)", email).First(&Usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &Usuario, nil
}

// ExistsByEmailExcept verifica se existe um usuário com o email especificado, exceto o usuário com o ID especificado
func (r *GormUsuarioRepository) VerificaOutroUsuarioPorEmail(email string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.Usuario{}).Where("LOWER(email) = LOWER(?) AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}

// Create cria um novo usuário
func (r *GormUsuarioRepository) Create(Usuario *models.Usuario) error {
	return r.GetDB().Create(Usuario).Error
}

// Update atualiza um usuário existente
func (r *GormUsuarioRepository) Update(Usuario *models.Usuario) error {
	return r.GetDB().Save(Usuario).Error
}

// Delete exclui um usuário (soft delete)
func (r *GormUsuarioRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.Usuario{}, id).Error
}
