package routes

import (
	"github.com/gin-gonic/gin"

	"msd-service/server/internal/api/handlers"
	"msd-service/server/internal/middleware"
)

// SetupComportamentoAlvoRoutes configura as rotas relacionadas a comportamentos alvo
func SetupComportamentoAlvoRoutes(router *gin.RouterGroup, handler *handlers.ComportamentoAlvoHandler, authMiddleware middleware.AuthMiddleware) {
	comportamentos := router.Group("/comportamentos")
	comportamentos.Use(authMiddleware.RequireAuth())
	{
		comportamentos.POST("", handler.CreateComportamento)
		comportamentos.GET("", handler.ListComportamentos)
		comportamentos.GET("/:id", handler.GetComportamento)
		comportamentos.PUT("/:id", handler.UpdateComportamento)
		comportamentos.DELETE("/:id", handler.DeleteComportamento)
	}

	// Rotas aninhadas para comportamentos de um paciente específico
	pacientes := router.Group("/pacientes")
	pacientes.Use(authMiddleware.RequireAuth())
	{
		pacientes.GET("/:paciente_id/comportamentos", handler.ListComportamentosByPaciente)
	}
}
