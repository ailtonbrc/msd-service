package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"msd-service/server/internal/api/handlers"
	"msd-service/server/internal/api/routes"
	"msd-service/server/internal/middleware"
	"msd-service/server/internal/repository"
	"msd-service/server/internal/service"
)

// Server representa o servidor HTTP da aplicação
type Server struct {
	router           *gin.Engine
	httpServer       *http.Server
	pacienteRepo     repository.PacienteRepository
	pacienteService  *service.PacienteService
	pacienteHandler  *handlers.PacienteHandler
	terapiaRepo      repository.TerapiaRepository
	terapiaService   *service.TerapiaService
	terapiaHandler   *handlers.TerapiaHandler
	sessaoRepo       repository.SessaoRepository
	sessaoService    *service.SessaoService
	sessaoHandler    *handlers.SessaoHandler
	objetivoRepo     repository.ObjetivoTerapeuticoRepository
	objetivoService  *service.ObjetivoTerapeuticoService
	objetivoHandler  *handlers.ObjetivoTerapeuticoHandler
	programaRepo     repository.ProgramaABARepository
	programaService  *service.ProgramaABAService
	programaHandler  *handlers.ProgramaABAHandler
	comportamentoRepo repository.ComportamentoAlvoRepository
	comportamentoService *service.ComportamentoAlvoService
	comportamentoHandler *handlers.ComportamentoAlvoHandler
	authMiddleware   middleware.AuthMiddleware
}

// NewServer cria uma nova instância do servidor
func NewServer(db *gorm.DB, jwtSecret string) *Server {
	router := gin.Default()

	// Configuração do CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Configuração do Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Inicialização dos componentes
	authMiddleware := middleware.NewJWTAuthMiddleware(jwtSecret)
	
	// Repositórios
	pacienteRepo := repository.NewGormPacienteRepository(db)
	terapiaRepo := repository.NewGormTerapiaRepository(db)
	sessaoRepo := repository.NewGormSessaoRepository(db)
	objetivoRepo := repository.NewGormObjetivoTerapeuticoRepository(db)
	programaRepo := repository.NewGormProgramaABARepository(db)
	comportamentoRepo := repository.NewGormComportamentoAlvoRepository(db)
	
	// Serviços
	pacienteService := service.NewPacienteService(pacienteRepo)
	terapiaService := service.NewTerapiaService(terapiaRepo)
	sessaoService := service.NewSessaoService(sessaoRepo)
	objetivoService := service.NewObjetivoTerapeuticoService(objetivoRepo)
	programaService := service.NewProgramaABAService(programaRepo)
	comportamentoService := service.NewComportamentoAlvoService(comportamentoRepo)
	
	// Handlers
	pacienteHandler := handlers.NewPacienteHandler(pacienteService)
	terapiaHandler := handlers.NewTerapiaHandler(terapiaService)
	sessaoHandler := handlers.NewSessaoHandler(sessaoService)
	objetivoHandler := handlers.NewObjetivoTerapeuticoHandler(objetivoService)
	programaHandler := handlers.NewProgramaABAHandler(programaService)
	comportamentoHandler := handlers.NewComportamentoAlvoHandler(comportamentoService)

	server := &Server{
		router:           router,
		pacienteRepo:     pacienteRepo,
		pacienteService:  pacienteService,
		pacienteHandler:  pacienteHandler,
		terapiaRepo:      terapiaRepo,
		terapiaService:   terapiaService,
		terapiaHandler:   terapiaHandler,
		sessaoRepo:       sessaoRepo,
		sessaoService:    sessaoService,
		sessaoHandler:    sessaoHandler,
		objetivoRepo:     objetivoRepo,
		objetivoService:  objetivoService,
		objetivoHandler:  objetivoHandler,
		programaRepo:     programaRepo,
		programaService:  programaService,
		programaHandler:  programaHandler,
		comportamentoRepo: comportamentoRepo,
		comportamentoService: comportamentoService,
		comportamentoHandler: comportamentoHandler,
		authMiddleware:   authMiddleware,
		httpServer: &http.Server{
			Addr:    ":" + os.Getenv("PORT"),
			Handler: router,
		},
	}

	server.setupRoutes()
	return server
}

// setupRoutes configura as rotas da API
func (s *Server) setupRoutes() {
	// Rota de saúde
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Grupo de rotas da API v1
	v1 := s.router.Group("/api/v1")
	routes.SetupPacienteRoutes(v1, s.pacienteHandler, s.authMiddleware)
	routes.SetupTerapiaRoutes(v1, s.terapiaHandler, s.authMiddleware)
	routes.SetupSessaoRoutes(v1, s.sessaoHandler, s.authMiddleware)
	routes.SetupObjetivoTerapeuticoRoutes(v1, s.objetivoHandler, s.authMiddleware)
	routes.SetupProgramaABARoutes(v1, s.programaHandler, s.authMiddleware)
	routes.SetupComportamentoAlvoRoutes(v1, s.comportamentoHandler, s.authMiddleware)
}

// Start inicia o servidor HTTP
func (s *Server) Start() {
	// Iniciar o servidor em uma goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	// Configurar canal para capturar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Desligando o servidor...")

	// Contexto com timeout para o shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar o servidor: %v", err)
	}

	log.Println("Servidor desligado com sucesso")
}
