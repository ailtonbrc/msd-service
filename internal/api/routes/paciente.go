package routes

import (
	"clinica_server/config"
	"clinica_server/internal/api/handlers"
	"clinica_server/internal/api/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ConfiguraRotasDePaciente(router *gin.RouterGroup, db *gorm.DB) {
	pacienteHandler := handlers.NewPacienteHandler(db)

	// Obter configuração para middleware de autenticação
	cfg, _ := config.Load()

	// Grupo de rotas de usuários (todas protegidas)
	pacientes := router.Group("/pacientes")
	pacientes.Use(middlewares.AuthMiddleware(cfg))
	{
		// TODO Verificar como vai tratar a questão dos Perfis de Acesso
		// users.GET("", middlewares.RequirePermission("users.view"), userHandler.GetUsers)
		// users.GET("/:id", middlewares.RequirePermission("users.view"), userHandler.GetUser)
		// users.POST("", middlewares.RequirePermission("users.create"), userHandler.CreateUser)
		// users.PUT("/:id", middlewares.RequirePermission("users.edit"), userHandler.UpdateUser)
		// users.DELETE("/:id", middlewares.RequirePermission("users.delete"), userHandler.DeleteUser)
		// users.PUT("/:id/password", userHandler.ChangePassword) // Permissão verificada no handler

		pacientes.GET("", pacienteHandler.BuscaPacientes)
		pacientes.GET("/:id", pacienteHandler.BuscaPaciente)
		pacientes.POST("", pacienteHandler.CreatePaciente)
		pacientes.PUT("/:id", pacienteHandler.UpdatePaciente)
		pacientes.DELETE("/:id", pacienteHandler.DeletePaciente)
	}
}