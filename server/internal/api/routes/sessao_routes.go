package routes

import (
	"github.com/gin-gonic/gin"

	"msd-service/server/internal/api/handlers"
	"msd-service/server/internal/middleware"
)

// SetupSessaoRoutes configura as rotas relacionadas a sessões
func SetupSessaoRoutes(router *gin.RouterGroup, handler *handlers.SessaoHandler, authMiddleware middleware.AuthMiddleware) {
	sessoes := router.Group("/sessoes")
	sessoes.Use(authMiddleware.RequireAuth())
	{
		sessoes.POST("", handler.CreateSessao)
		sessoes.GET("", handler.ListSessoes)
		sessoes.GET("/:id", handler.GetSessao)
		sessoes.PUT("/:id", handler.UpdateSessao)
		sessoes.DELETE("/:id", handler.DeleteSessao)
	}

	// Rotas aninhadas para sessões de um paciente específico
	pacientes := router.Group("/pacientes")
	pacientes.Use(authMiddleware.RequireAuth())
	{
		pacientes.GET("/:paciente_id/sessoes", handler.ListSessoesByPaciente)
	}
}
