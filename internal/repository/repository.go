package repository

import (
	"context"

	"gorm.io/gorm"
)

// Repository é a interface base para todos os repositories
type Repository interface {
	GetDB() *gorm.DB
	WithContext(ctx context.Context) Repository
	WithTx(tx *gorm.DB) Repository
}

// BaseRepository implementa a interface Repository
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository cria um novo repository base
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB retorna a instância do banco de dados
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}

// WithContext retorna um novo repository com o contexto especificado
func (r *BaseRepository) WithContext(ctx context.Context) Repository {
	return &BaseRepository{db: r.db.WithContext(ctx)}
}

// WithTx retorna um novo repository com a transação especificada
func (r *BaseRepository) WithTx(tx *gorm.DB) Repository {
	return &BaseRepository{db: tx}
}
