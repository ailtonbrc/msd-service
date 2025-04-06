package routes

import (
	"clinica_server/config"
	"clinica_server/internal/api/handlers"
	"clinica_server/internal/api/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfiguraRotasDeAutenticacao configura as rotas de autenticação
func ConfiguraRotasDeAutenticacao(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(db, cfg)

	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh-token", authHandler.RefreshToken)

		// Rotas protegidas
		protected := auth.Group("")
		protected.Use(middlewares.AuthMiddleware(cfg))
		{
			protected.POST("/logout", authHandler.Logout)
			protected.GET("/me", authHandler.GetMe)
		}
	}
}
