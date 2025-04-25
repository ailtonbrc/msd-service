package routes

import (
	"github.com/gin-gonic/gin"

	"msd-service/server/internal/api/handlers"
	"msd-service/server/internal/middleware"
)

// SetupObjetivoTerapeuticoRoutes configura as rotas relacionadas a objetivos terapêuticos
func SetupObjetivoTerapeuticoRoutes(router *gin.RouterGroup, handler *handlers.ObjetivoTerapeuticoHandler, authMiddleware middleware.AuthMiddleware) {
	objetivos := router.Group("/objetivos")
	objetivos.Use(authMiddleware.RequireAuth())
	{
		objetivos.POST("", handler.CreateObjetivo)
		objetivos.GET("", handler.ListObjetivos)
		objetivos.GET("/:id", handler.GetObjetivo)
		objetivos.PUT("/:id", handler.UpdateObjetivo)
		objetivos.DELETE("/:id", handler.DeleteObjetivo)
	}

	// Rotas aninhadas para objetivos de um paciente específico
	pacientes := router.Group("/pacientes")
	pacientes.Use(authMiddleware.RequireAuth())
	{
		pacientes.GET("/:paciente_id/objetivos", handler.ListObjetivosByPaciente)
	}
}
