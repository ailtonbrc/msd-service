package seeders

import "gorm.io/gorm"

func RunAll(db *gorm.DB) {
	SeedUsuarioAdm(db) // Cria o usuário administrador
}
