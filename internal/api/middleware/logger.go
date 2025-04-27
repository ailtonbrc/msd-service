package middleware

import (
	"time"

	"clinica_server/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoggerMiddleware registra informações sobre as requisições
func LoggerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tempo de início
		startTime := time.Now()

		// Processar requisição
		c.Next()

		// Tempo de término
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Obter informações do usuário
		userID, exists := c.Get("userID")
		if !exists {
			userID = 0
		}

		// Registrar log no banco de dados para operações importantes
		if c.Request.Method != "GET" && c.Writer.Status() < 400 {
			log := models.SystemLog{
				UserID:    userID.(*uint),
				Action:    c.Request.Method + " " + c.Request.URL.Path,
				IPAddress: c.ClientIP(),
				Details: map[string]interface{}{
					"status":     c.Writer.Status(),
					"latency_ms": latency.Milliseconds(),
					"user_agent": c.Request.UserAgent(),
				},
			}

			// Extrair informações sobre a entidade afetada
			if entityType := c.Param("entity_type"); entityType != "" {
				log.EntityType = entityType
			}

			if entityID := c.Param("id"); entityID != "" {
				log.EntityID = entityID
			}

			// Salvar log de forma assíncrona
			go func(log models.SystemLog) {
				db.Create(&log)
			}(log)
		}
	}
}
