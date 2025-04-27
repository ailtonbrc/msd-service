package utils

import (
	"fmt"
	"time"

	"clinica_server/config"
)

// DisplayServerInfo exibe informações amigáveis sobre o servidor
func DisplayServerInfo(cfg *config.Config) {
	fmt.Println("\n========================================================")
	fmt.Println("           CLÍNICA TEA - API SERVER                    ")
	fmt.Println("========================================================")
	fmt.Printf("✅ Servidor iniciado com sucesso na porta: %s\n", cfg.Server.Port)
	fmt.Printf("🌐 Ambiente: %s\n", cfg.Environment)
	fmt.Printf("🔐 JWT configurado com expiração: %v\n", cfg.JWT.AccessTokenExp)
	fmt.Println("📚 API Endpoints disponíveis:")
	fmt.Println("   - Autenticação: /api/auth")
	fmt.Println("   - Usuários:     /api/usuarios")
	fmt.Println("   - Pacientes:    /api/pacientes")
	fmt.Println("🔍 Verificação de saúde: /health")
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("⏱️  Iniciado em: %s\n", time.Now().Format("02/01/2006 15:04:05"))
	fmt.Println("========================================================\n")
}