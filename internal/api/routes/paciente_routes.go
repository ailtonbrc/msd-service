package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"clinica_server/internal/api/handlers"
	"clinica_server/internal/api/middleware"
)

// SetupPacienteRoutes configura as rotas para o módulo de pacientes
// Recebe o router, o handler, o middleware de autenticação e o logger
func SetupPacienteRoutes(router *gin.RouterGroup, handler *handlers.PacienteHandler, authMiddleware *middleware.AuthMiddleware, logger *zap.Logger) {
	// Grupo de rotas para pacientes
	pacienteGroup := router.Group("/pacientes")
	
	// Aplicar middleware de autenticação em todas as rotas
	pacienteGroup.Use(authMiddleware.AuthRequired())
	
	// Rotas GET (leitura)
	pacienteGroup.GET("", authMiddleware.CheckPermission("pacientes:read"), handler.GetAll)
	pacienteGroup.GET("/:id", authMiddleware.CheckPermission("pacientes:read"), handler.GetByID)
	pacienteGroup.GET("/busca", authMiddleware.CheckPermission("pacientes:read"), handler.Search)
	
	// Rotas POST (criação)
	pacienteGroup.POST("", authMiddleware.CheckPermission("pacientes:create"), handler.Create)
	
	// Rotas PUT (atualização)
	pacienteGroup.PUT("/:id", authMiddleware.CheckPermission("pacientes:update"), handler.Update)
	
	// Rotas DELETE (exclusão)
	pacienteGroup.DELETE("/:id", authMiddleware.CheckPermission("pacientes:delete"), handler.Delete)
	
	logger.Info("Rotas de pacientes configuradas")
}