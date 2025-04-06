package seeders

import (
	"log"

	"clinica_server/internal/models"

	"gorm.io/gorm"

	"clinica_server/internal/utils"
)

func SeedUsuarioAdm(db *gorm.DB) {
	// Hash da senha
	passwordHash, err := utils.HashPassword("123456")
	if err != nil {
		log.Fatal("Erro ao configurar a senha de adm:", err)
	}

	// Criar usuário administrador
	adminUser := models.Usuario{
		Nome:   "Administrador",
		Email:  "admin@sistema.com",
		Senha:  passwordHash,
		Perfil: "ADMIN",
		Ativo:  true,
	}

	// Evitar duplicação
	if err := db.Where("email = ?", adminUser.Email).FirstOrCreate(&adminUser).Error; err != nil {
		log.Printf("Erro ao criar usuário ADMIN: %v", err)
	} else {
		log.Println("Usuário ADMIN criado ou já existente!")
	}
}
