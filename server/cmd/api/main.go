package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"msd-service/server/internal/api/server"
	"msd-service/server/internal/models"
)

// @title MSD Service API
// @version 1.0
// @description API para gerenciamento de clínica de terapias TEA/ABA
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite 'Bearer ' seguido do token JWT
func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Configurar conexão com o banco de dados
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Configurar logger do GORM
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// Conectar ao banco de dados
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}

	// Executar migrações
	if err := db.AutoMigrate(
		&models.Paciente{},
		&models.Terapia{},
		&models.Sessao{},
		&models.ObjetivoTerapeutico{},
		&models.ProgressoObjetivo{},
		&models.ProgramaABA{},
		&models.EtapaPrograma{},
		&models.TipoPrompt{},
		&models.ColetaABA{},
		&models.ComportamentoAlvo{},
		&models.RegistroComportamento{},
	); err != nil {
		log.Fatalf("Falha ao executar migrações: %v", err)
	}

	// Obter chave secreta para JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Variável de ambiente JWT_SECRET não definida")
	}

	// Definir porta padrão se não estiver configurada
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}

	// Iniciar o servidor
	srv := server.NewServer(db, jwtSecret)
	log.Printf("Servidor iniciado na porta %s", os.Getenv("PORT"))
	srv.Start()
}
