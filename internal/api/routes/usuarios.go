package routes

import (
	"clinica_server/config"
	"clinica_server/internal/api/handlers"
	"clinica_server/internal/api/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfiguraRotasDeUsuario configura as rotas de usuários
func ConfiguraRotasDeUsuario(router *gin.RouterGroup, db *gorm.DB) {
	userHandler := handlers.NewUsuarioHandler(db)

	// Obter configuração para middleware de autenticação
	cfg, _ := config.Load()

	// Grupo de rotas de usuários (todas protegidas)
	users := router.Group("/users")
	users.Use(middlewares.AuthMiddleware(cfg))
	{
		// TODO Verificar como vai tratar a questão dos Perfis de Acesso
		// users.GET("", middlewares.RequirePermission("users.view"), userHandler.GetUsers)
		// users.GET("/:id", middlewares.RequirePermission("users.view"), userHandler.GetUser)
		// users.POST("", middlewares.RequirePermission("users.create"), userHandler.CreateUser)
		// users.PUT("/:id", middlewares.RequirePermission("users.edit"), userHandler.UpdateUser)
		// users.DELETE("/:id", middlewares.RequirePermission("users.delete"), userHandler.DeleteUser)
		// users.PUT("/:id/password", userHandler.ChangePassword) // Permissão verificada no handler

		users.GET("", userHandler.BuscaUsuarios)
		users.GET("/:id", userHandler.BuscaUsuario)
		users.POST("", userHandler.CreateUsuario)
		users.PUT("/:id", userHandler.UpdateUsuario)
		users.DELETE("/:id", userHandler.DeleteUsuario)
		users.PUT("/:id/password", userHandler.ChangePassword) // Permissão verificada no handler
	}
}
