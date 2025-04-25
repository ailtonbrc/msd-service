package routes

import (
	"github.com/gin-gonic/gin"

	"msd-service/server/internal/api/handlers"
	"msd-service/server/internal/middleware"
)

// SetupProgramaABARoutes configura as rotas relacionadas a programas ABA
func SetupProgramaABARoutes(router *gin.RouterGroup, handler *handlers.ProgramaABAHandler, authMiddleware middleware.AuthMiddleware) {
	programas := router.Group("/programas")
	programas.Use(authMiddleware.RequireAuth())
	{
		programas.POST("", handler.CreatePrograma)
		programas.GET("", handler.ListProgramas)
		programas.GET("/:id", handler.GetPrograma)
		programas.PUT("/:id", handler.UpdatePrograma)
		programas.DELETE("/:id", handler.DeletePrograma)
	}

	// Rotas aninhadas para programas de um paciente específico
	pacientes := router.Group("/pacientes")
	pacientes.Use(authMiddleware.RequireAuth())
	{
		pacientes.GET("/:paciente_id/programas", handler.ListProgramasByPaciente)
	}
}
