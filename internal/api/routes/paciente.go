package routes

import (
	"clinica_server/config"
	"clinica_server/internal/api/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPacienteRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
    handler := handlers.NewPacienteHandler(db)
    pacientes := rg.Group("/pacientes")
   // pacientes.Use(middlewares.AuthMiddleware(cfg))
    {
        pacientes.GET("", handler.GetPacientes)
        pacientes.GET("/:id", handler.GetPaciente)
        pacientes.POST("", handler.CreatePaciente)
        pacientes.PUT("/:id", handler.UpdatePaciente)
        pacientes.DELETE("/:id", handler.DeletePaciente)
    }
}
