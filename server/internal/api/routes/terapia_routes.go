package routes

import (
	"github.com/gin-gonic/gin"

	"msd-service/server/internal/api/handlers"
	"msd-service/server/internal/middleware"
)

// SetupTerapiaRoutes configura as rotas relacionadas a terapias
func SetupTerapiaRoutes(router *gin.RouterGroup, handler *handlers.TerapiaHandler, authMiddleware middleware.AuthMiddleware) {
	terapias := router.Group("/terapias")
	terapias.Use(authMiddleware.RequireAuth())
	{
		terapias.POST("", handler.CreateTerapia)
		terapias.GET("", handler.ListTerapias)
		terapias.GET("/:id", handler.GetTerapia)
		terapias.PUT("/:id", handler.UpdateTerapia)
		terapias.DELETE("/:id", handler.DeleteTerapia)
	}
}
