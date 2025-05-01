package migrations

import (
	"clinica_server/internal/models"
	"log"

	"gorm.io/gorm"
)

// MigrateDB executa as migrações do banco de dados
func MigrateDB(db *gorm.DB) error {
	log.Println("Iniciando migrações do banco de dados...")

	// Lista de todos os modelos para migração
	models := []interface{}{
		&models.Usuario{},
		&models.Paciente{},
		&models.Clinica{},
		&models.SystemLog{},
	}

	// Executar migrações
	err := db.AutoMigrate(models...)
	if err != nil {
		log.Printf("Erro ao executar migrações: %v", err)
		return err
	}

	log.Println("Migrações concluídas com sucesso!")
	return nil
}

// Adicionar ao arquivo migrate.go ou criar se não existir

// // AutoMigrate realiza a migração automática das tabelas
// func AutoMigrate(db *gorm.DB) error {
// 	// Migrar tabelas
// 	err := db.AutoMigrate(
// 		&Usuario{},
// 		&Paciente{}, // Adicionar modelo de paciente
// 		// Outros modelos...
// 	)
	
// 	return err
// }
