package routes

import (
    "clinica_server/internal/api/handlers"
    "clinica_server/internal/api/middleware"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// ConfiguraRotasDeUsuario configura as rotas de usuários
func ConfiguraRotasDeUsuario(router *gin.RouterGroup, db *gorm.DB, authMiddleware *middleware.AuthMiddleware) {
    userHandler := handlers.NewUsuarioHandler(db)

    // Grupo de rotas de usuários (todas protegidas)
    usuarios := router.Group("/usuarios")
    usuarios.Use(authMiddleware.RequireAuth())
    {
        // TODO Verificar como vai tratar a questão dos Perfis de Acesso
        // users.GET("", middleware.RequirePermission("users.view"), userHandler.GetUsers)
        // users.GET("/:id", middleware.RequirePermission("users.view"), userHandler.GetUser)
        // users.POST("", middleware.RequirePermission("users.create"), userHandler.CreateUser)
        // users.PUT("/:id", middleware.RequirePermission("users.edit"), userHandler.UpdateUser)
        // users.DELETE("/:id", middleware.RequirePermission("users.delete"), userHandler.DeleteUser)
        // users.PUT("/:id/password", userHandler.ChangePassword) // Permissão verificada no handler

        usuarios.GET("", userHandler.BuscaUsuarios)
        usuarios.GET("/:id", userHandler.BuscaUsuario)
        usuarios.POST("", userHandler.CreateUsuario)
        usuarios.PUT("/:id", userHandler.UpdateUsuario)
        usuarios.DELETE("/:id", userHandler.DeleteUsuario)
        usuarios.PUT("/:id/password", userHandler.ChangePassword) // Permissão verificada no handler
    }
}