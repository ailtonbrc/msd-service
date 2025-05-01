package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clinica_server/config"
	"clinica_server/internal/api/handlers"
	"clinica_server/internal/api/middleware"
	"clinica_server/internal/api/routes"
	"clinica_server/internal/auth"
	"clinica_server/internal/models"
	"clinica_server/internal/repository"
	"clinica_server/internal/service"
	"clinica_server/internal/utils"
	"clinica_server/internal/validator"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Server representa o servidor HTTP
type Server struct {
	router *gin.Engine
	cfg    *config.Config
	db     *gorm.DB
	logger *zap.Logger
}

// NewServer cria uma nova instância do servidor
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	// Configurar logger
	logger, err := zap.NewProduction()
	if cfg.Environment == "development" {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("Erro ao configurar logger: %v", err)
	}

	// Configurar o Gin
	configureGin(cfg.Environment)

	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Server{
		router: router,
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

// configureGin configura o framework Gin
func configureGin(environment string) {
	// Configurar o modo do Gin
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// Redirecionar logs para um arquivo
	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
}

// Run inicia o servidor HTTP
func (s *Server) Run() error {
	// Configurar rotas
	s.setupRoutes()

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	// Canal para capturar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em uma goroutine
	go func() {
		s.logger.Info("Servidor iniciado", zap.String("porta", s.cfg.Server.Port))

		// Exibir informações do servidor
		utils.DisplayServerInfo(s.cfg)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	s.logger.Info("Desligando servidor...")

	// Contexto com timeout para desligamento gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Desligar servidor
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Fatal("Erro ao desligar servidor", zap.Error(err))
	}

	s.logger.Info("Servidor desligado com sucesso")
	return nil
}

// setupRoutes configura todas as rotas da API
func (s *Server) setupRoutes() {
	// Grupo de rotas da API
	api := s.router.Group("/api")

	// Configurar serviço JWT
	jwtService := auth.NewJWTService(s.cfg.JWT.Secret)
	
	// Configurar middleware de autenticação
	authMiddleware := middleware.NewAuthMiddleware(jwtService, s.logger)
	
	// Configurar repositórios
	usuarioRepo := repository.NewUsuarioRepository(s.db)
	pacienteRepo := repository.NewPacienteRepository(s.db)
	
	// Configurar conversores de DTO
	pacienteDTOConverter := models.NewPacienteDTOConverter()
	
	// Configurar validadores
	pacienteValidator := validator.NewPacienteValidator(pacienteRepo)
	
	// Configurar serviços
	authService := service.NewAuthService(usuarioRepo, jwtService, s.cfg.JWT.AccessTokenExp, s.cfg.JWT.RefreshTokenExp)
	pacienteService := service.NewPacienteService(pacienteRepo, pacienteValidator, pacienteDTOConverter)
	
	// Configurar handlers
	authHandler := handlers.NewAuthHandler(authService, s.logger)
	usuarioHandler := handlers.NewUsuarioHandler(usuarioRepo, s.logger)
	pacienteHandler := handlers.NewPacienteHandler(pacienteService, s.logger)

	// Configurar rotas para cada módulo
	routes.ConfiguraRotasDeAutenticacao(api, s.db, authMiddleware)
	routes.ConfiguraRotasDeUsuario(api, s.db, authMiddleware)
	
	// Configurar rotas de pacientes
	routes.SetupPacienteRoutes(api, pacienteHandler, authMiddleware, s.logger)
	
	// Rota de verificação de saúde do servidor
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
			"version": "1.0.0", // Versão da API
		})
	})
	
	// Rota para endpoints não encontrados
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Não encontrado",
			"message": "Endpoint não encontrado",
		})
	})
	
	s.logger.Info("Rotas configuradas com sucesso")
}

// Adicionar ao método setupRoutes no arquivo server.go

// Configurar repositórios
usuarioRepo := repository.NewUsuarioRepository(s.db)
pacienteRepo := repository.NewPacienteRepository(s.db)

// Configurar conversores de DTO
pacienteDTOConverter := models.NewPacienteDTOConverter()

// Configurar validadores
pacienteValidator := validator.NewPacienteValidator(pacienteRepo)

// Configurar serviços
authService := service.NewAuthService(usuarioRepo, jwtService, s.cfg.JWT.AccessTokenExp, s.cfg.JWT.RefreshTokenExp)
pacienteService := service.NewPacienteService(pacienteRepo, pacienteValidator, pacienteDTOConverter)

// Configurar handlers
authHandler := handlers.NewAuthHandler(authService, s.logger)
usuarioHandler := handlers.NewUsuarioHandler(usuarioRepo, s.logger)
pacienteHandler := handlers.NewPacienteHandler(pacienteService, s.logger)

// Configurar rotas para cada módulo
routes.ConfiguraRotasDeAutenticacao(api, s.db, authMiddleware)
routes.ConfiguraRotasDeUsuario(api, s.db, authMiddleware)

// Configurar rotas de pacientes
routes.SetupPacienteRoutes(api, pacienteHandler, authMiddleware, s.logger)