package main

import (
	"log"

	"clinica_server/config"
	"clinica_server/internal/api/server"
	"clinica_server/internal/db"

	"go.uber.org/zap"
)
func main() {
	// Carregar configurações clinica_server
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Configurar logger
	logger, err := zap.NewProduction()
	if cfg.Environment == "development" {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("Erro ao configurar logger: %v", err)
	}
	defer logger.Sync()
	
	// Inicializar banco de dados
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Inicializar e executar o servidor
	s := server.NewServer(cfg, database)
	
	if err := s.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
