package seeders

import "gorm.io/gorm"

func RunAll(db *gorm.DB) {
	SeedUsuarioAdm(db) // Cria o usuário administrador
}

// SeedPermissoes adiciona as permissões padrão ao banco de dados
func SeedPermissoes(db *gorm.DB) error {
	// Permissões para pacientes
	permissoes := []Permissao{
		{Nome: "pacientes:read", Descricao: "Visualizar pacientes"},
		{Nome: "pacientes:create", Descricao: "Criar pacientes"},
		{Nome: "pacientes:update", Descricao: "Atualizar pacientes"},
		{Nome: "pacientes:delete", Descricao: "Excluir pacientes"},
		// Outras permissões...
	}
	
	// Inserir permissões
	for _, p := range permissoes {
		var count int64
		db.Model(&Permissao{}).Where("nome = ?", p.Nome).Count(&count)
		if count == 0 {
			db.Create(&p)
		}
	}
	
	return nil
}