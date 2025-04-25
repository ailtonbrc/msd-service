package routes

import (
	"github.com/gin-gonic/gin"

	"msd-service/server/internal/api/handlers"
	"msd-service/server/internal/middleware"
)

// SetupPacienteRoutes configura as rotas relacionadas a pacientes
func SetupPacienteRoutes(router *gin.RouterGroup, handler *handlers.PacienteHandler, authMiddleware middleware.AuthMiddleware) {
	pacientes := router.Group("/pacientes")
	pacientes.Use(authMiddleware.RequireAuth())
	{
		pacientes.POST("", handler.CreatePaciente)
		pacientes.GET("", handler.ListPacientes)
		pacientes.GET("/:id", handler.GetPaciente)
		pacientes.PUT("/:id", handler.UpdatePaciente)
		pacientes.DELETE("/:id", handler.DeletePaciente)
	}
}
