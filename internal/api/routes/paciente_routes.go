// internal/api/routes/paciente_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"clinica_server/internal/api/handlers"
	middleware "clinica_server/internal/api/middleware"
)

// SetupPacienteRoutes configura as rotas relacionadas a pacientes
func SetupPacienteRoutes(
	router *gin.Engine,
	pacienteHandler *handlers.PacienteHandler,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) {
	// Grupo de rotas para API
	api := router.Group("/api")
	
	// Grupo de rotas para pacientes
	pacientes := api.Group("/pacientes")
	
	// Aplicar middleware de autenticação em todas as rotas de pacientes
	pacientes.Use(authMiddleware.RequireAuth())
	
	// Rotas que exigem permissão de visualização (pacientes:view)
	pacientes.GET("", middleware.RequirePermission("pacientes:view"), pacienteHandler.List)
	pacientes.GET("/:id", middleware.RequirePermission("pacientes:view"), pacienteHandler.GetByID)
	pacientes.GET("/search", middleware.RequirePermission("pacientes:view"), pacienteHandler.Search)
	pacientes.GET("/cpf/:cpf", middleware.RequirePermission("pacientes:view"), pacienteHandler.GetByCPF)
	pacientes.GET("/:id/idade", middleware.RequirePermission("pacientes:view"), pacienteHandler.CalcularIdade)
	
	// Rotas que exigem permissão de criação (pacientes:create)
	pacientes.POST("", middleware.RequirePermission("pacientes:create"), pacienteHandler.Create)
	
	// Rotas que exigem permissão de atualização (pacientes:update)
	pacientes.PUT("/:id", middleware.RequirePermission("pacientes:update"), pacienteHandler.Update)
	pacientes.PATCH("/:id/diagnostico", middleware.RequirePermission("pacientes:update"), pacienteHandler.AtualizarDiagnostico)
	
	// Rotas que exigem permissão de exclusão (pacientes:delete)
	pacientes.DELETE("/:id", middleware.RequirePermission("pacientes:delete"), pacienteHandler.Delete)
	
	// Registrar as rotas no logger
	logger.Info("Rotas de pacientes configuradas")
}

// RegisterPacienteRoutes registra as rotas de pacientes no router principal
func RegisterPacienteRoutes(
	router *gin.Engine,
	pacienteHandler *handlers.PacienteHandler,
	authMiddleware *middleware.AuthMiddleware,
	logger *zap.Logger,
) {
	SetupPacienteRoutes(router, pacienteHandler, authMiddleware, logger)
}