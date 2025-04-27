// internal/api/routes/auth.go
package routes

import (
	"clinica_server/internal/api/handlers"
	"clinica_server/internal/api/middleware"
	"clinica_server/internal/repository"
	"clinica_server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfiguraRotasDeAutenticacao configura as rotas de autenticação
func ConfiguraRotasDeAutenticacao(router *gin.RouterGroup, db *gorm.DB, authMiddleware *middleware.AuthMiddleware) {
	// Criar repositório de usuários
	userRepo := repository.NewUsuarioRepository(db)
	
	// Criar serviço de autenticação usando o repositório e o JWTService
	authService := service.NewAuthService(userRepo, authMiddleware.GetJWTService())
	
	// Criar handler de autenticação
	authHandler := handlers.NewAuthHandler(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh-token", authHandler.RefreshToken)

		// Rotas protegidas
		protected := auth.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.POST("/logout", authHandler.Logout)
			protected.GET("/me", authHandler.GetMe)
		}
	}
}