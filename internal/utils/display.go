package utils

import (
	"fmt"
	"time"

	"clinica_server/config"
)

// DisplayServerInfo exibe informações amigáveis sobre o servidor
func DisplayServerInfo(cfg *config.Config) {
	// Adicionar um log de depuração
	fmt.Printf("DEBUG: JWT.AccessTokenExp = %v (%d minutos)\n", cfg.JWT.AccessTokenExp,int(cfg.JWT.AccessTokenExp.Minutes()))

	fmt.Println("\n========================================================")
	fmt.Println("           CLÍNICA TEA - API SERVER                    ")
	fmt.Println("========================================================")
	fmt.Printf("✅ Servidor iniciado com sucesso na porta: %s\n", cfg.Server.Port)
	fmt.Printf("🌐 Ambiente: %s\n", cfg.Environment)
	
	// Formatar o tempo de expiração de forma mais amigável
	expiracaoFormatada := formatarDuracao(cfg.JWT.AccessTokenExp)
	fmt.Printf("🔐 JWT configurado com expiração: %s\n", expiracaoFormatada)
	
	fmt.Println("📚 API Endpoints disponíveis:")
	fmt.Println("   - Autenticação: /api/auth")
	fmt.Println("   - Usuários:     /api/usuarios")
	fmt.Println("   - Pacientes:    /api/pacientes")
	fmt.Println("🔍 Verificação de saúde: /health")
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("⏱️  Iniciado em: %s\n", time.Now().Format("02/01/2006 15:04:05"))
	fmt.Println("========================================================\n")
}

// formatarDuracao converte uma duração em uma string amigável
func formatarDuracao(d time.Duration) string {
	minutos := int(d.Minutes())
	
	if minutos < 60 {
		return fmt.Sprintf("%d minutos", minutos)
	} else if minutos < 1440 {
		horas := minutos / 60
		minutosRestantes := minutos % 60
		if minutosRestantes == 0 {
			return fmt.Sprintf("%d horas", horas)
		} else {
			return fmt.Sprintf("%d horas e %d minutos", horas, minutosRestantes)
		}
	} else {
		dias := minutos / 1440
		horasRestantes := (minutos % 1440) / 60
		if horasRestantes == 0 {
			return fmt.Sprintf("%d dias", dias)
		} else {
			return fmt.Sprintf("%d dias e %d horas", dias, horasRestantes)
		}
	}
}