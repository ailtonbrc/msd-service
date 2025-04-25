#!/bin/bash

# Script para gerar a estrutura do projeto Simple ERP Service

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Iniciando a criação do projeto Simple ERP Service...${NC}"

# Nome do módulo Go
MODULE_NAME="github.com/yourusername/simple-erp-service"

# Criar diretórios
echo -e "${GREEN}Criando estrutura de diretórios...${NC}"

mkdir -p cmd/api
mkdir -p config
mkdir -p internal/api/handlers
mkdir -p internal/api/middlewares
mkdir -p internal/api/routes
mkdir -p internal/models
mkdir -p internal/repository
mkdir -p internal/service
mkdir -p internal/utils
mkdir -p pkg/logger
mkdir -p pkg/validator
mkdir -p migrations

# Inicializar módulo Go
echo -e "${GREEN}Inicializando módulo Go...${NC}"
go mod init $MODULE_NAME

# Criar go.mod
cat > go.mod << 'EOF'
module github.com/yourusername/simple-erp-service

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.14.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.4
)

require (
	github.com/bytedance/sonic v1.10.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.15.5 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.5.0 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
EOF

# Criar arquivo .env
cat > .env << 'EOF'
# Configurações do Servidor
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10
SERVER_IDLE_TIMEOUT=60

# Configurações do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=erp_system
DB_SSLMODE=disable

# Configurações JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_EXP=15
JWT_REFRESH_EXP=10080
EOF

# Criar arquivo main.go
cat > cmd/api/main.go << 'EOF'
package main

import (
	"log"

	"github.com/yourusername/simple-erp-service/config"
	"github.com/yourusername/simple-erp-service/internal/api/server"
	"github.com/yourusername/simple-erp-service/internal/models"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Inicializar banco de dados
	db, err := models.InitDB(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Inicializar e executar o servidor
	s := server.NewServer(cfg, db)
	if err := s.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
EOF

# Criar arquivo config.go
cat > config/config.go << 'EOF'
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig armazena configurações do servidor HTTP
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig armazena configurações do banco de dados
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig armazena configurações do JWT
type JWTConfig struct {
	Secret           string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

// Load carrega as configurações do ambiente
func Load() (*Config, error) {
	// Carregar variáveis de ambiente do arquivo .env se existir
	_ = godotenv.Load()

	// Configurações do servidor
	port := getEnv("SERVER_PORT", "8080")
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("SERVER_IDLE_TIMEOUT", "60"))

	// Configurações do banco de dados
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "erp_system")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Configurações do JWT
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtAccessExp, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXP", "15"))    // 15 minutos
	jwtRefreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXP", "10080")) // 7 dias

	return &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		},
		Database: DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
		JWT: JWTConfig{
			Secret:           jwtSecret,
			AccessTokenExp:  time.Duration(jwtAccessExp) * time.Minute,
			RefreshTokenExp: time.Duration(jwtRefreshExp) * time.Minute,
		},
	}, nil
}

// DSN retorna a string de conexão com o banco de dados
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
EOF

# Criar arquivo server.go
cat > internal/api/server/server.go << 'EOF'
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/simple-erp-service/config"
	"github.com/yourusername/simple-erp-service/internal/api/routes"
	"gorm.io/gorm"
)

// Server representa o servidor HTTP
type Server struct {
	router *gin.Engine
	cfg    *config.Config
	db     *gorm.DB
}

// NewServer cria uma nova instância do servidor
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Server{
		router: router,
		cfg:    cfg,
		db:     db,
	}
}

// Run inicia o servidor HTTP
func (s *Server) Run() error {
	// Configurar rotas
	s.setupRoutes()

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	// Canal para capturar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em uma goroutine
	go func() {
		log.Printf("Servidor iniciado na porta %s", s.cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Desligando servidor...")

	// Contexto com timeout para desligamento gracioso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Desligar servidor
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}

	log.Println("Servidor desligado com sucesso")
	return nil
}

// setupRoutes configura todas as rotas da API
func (s *Server) setupRoutes() {
	// Grupo de rotas da API
	api := s.router.Group("/api")

	// Configurar rotas para cada módulo
	routes.SetupAuthRoutes(api, s.db, s.cfg)
	routes.SetupUserRoutes(api, s.db)
	routes.SetupRoleRoutes(api, s.db)
	routes.SetupProductRoutes(api, s.db)
	routes.SetupInventoryRoutes(api, s.db)
	routes.SetupCustomerRoutes(api, s.db)
	routes.SetupSupplierRoutes(api, s.db)
	routes.SetupSaleRoutes(api, s.db)
	routes.SetupPurchaseRoutes(api, s.db)
	routes.SetupFinancialRoutes(api, s.db)
	routes.SetupDashboardRoutes(api, s.db)
	routes.SetupSystemRoutes(api, s.db)
}
EOF

# Criar arquivo db.go
cat > internal/models/db.go << 'EOF'
package models

import (
	"github.com/yourusername/simple-erp-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB inicializa a conexão com o banco de dados
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Registrar modelos para auto-migração se necessário
	// Nota: Como já temos o esquema SQL, não precisamos de auto-migração
	// mas é bom ter isso configurado para desenvolvimento

	return db, nil
}
EOF

# Criar arquivo user.go
cat > internal/models/user.go << 'EOF'
package models

import (
	"time"

	"gorm.io/gorm"
)

// User representa um usuário do sistema
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;not null;unique" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Email        string         `gorm:"size:100;unique" json:"email"`
	Phone        string         `gorm:"size:20" json:"phone"`
	RoleID       uint           `json:"role_id"`
	Role         *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time     `json:"last_login"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserResponse é usado para retornar informações do usuário sem dados sensíveis
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	RoleID    uint      `json:"role_id"`
	Role      *Role     `json:"role,omitempty"`
	IsActive  bool      `json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converte um User para UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		RoleID:    u.RoleID,
		Role:      u.Role,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// CreateUserRequest representa os dados para criar um novo usuário
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

// UpdateUserRequest representa os dados para atualizar um usuário
type UpdateUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email" binding:"omitempty,email"`
	Phone  string `json:"phone"`
	RoleID uint   `json:"role_id"`
}

// ChangePasswordRequest representa os dados para alterar a senha
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}
EOF

# Criar arquivo role.go
cat > internal/models/role.go << 'EOF'
package models

import (
	"time"

	"gorm.io/gorm"
)

// Role representa um perfil de usuário
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null;unique" json:"name"`
	Description string         `json:"description"`
	Permissions []*Permission  `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateRoleRequest representa os dados para criar um novo perfil
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest representa os dados para atualizar um perfil
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateRolePermissionsRequest representa os dados para atualizar permissões de um perfil
type UpdateRolePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}
EOF

# Criar arquivo permission.go
cat > internal/models/permission.go << 'EOF'
package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission representa uma permissão no sistema
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;unique" json:"name"`
	Description string         `json:"description"`
	Module      string         `gorm:"size:50;not null" json:"module"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// PermissionsByModule agrupa permissões por módulo
type PermissionsByModule struct {
	Module      string       `json:"module"`
	Permissions []Permission `json:"permissions"`
}
EOF

# Criar arquivo system_log.go
cat > internal/models/system_log.go << 'EOF'
package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemLog representa um log do sistema
type SystemLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`
	Action     string         `gorm:"size:100;not null" json:"action"`
	EntityType string         `gorm:"size:50" json:"entity_type"`
	EntityID   string         `json:"entity_id"`
	Details    map[string]interface{} `gorm:"type:jsonb" json:"details"`
	IPAddress  string         `gorm:"size:45" json:"ip_address"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
EOF

# Criar arquivo jwt.go
cat > internal/utils/jwt.go << 'EOF'
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/simple-erp-service/config"
)

// JWTClaims representa os claims do token JWT
type JWTClaims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	RoleID   uint     `json:"role_id"`
	Role     string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// GenerateAccessToken gera um novo token JWT de acesso
func GenerateAccessToken(userID uint, username string, roleID uint, role string, permissions []string, cfg *config.Config) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		Role:     role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.AccessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "simple-erp-service",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// GenerateRefreshToken gera um novo token JWT de refresh
func GenerateRefreshToken(userID uint, username string, cfg *config.Config) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshTokenExp)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "simple-erp-service",
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ValidateToken valida um token JWT
func ValidateToken(tokenString string, cfg *config.Config) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
EOF

# Criar arquivo password.go
cat > internal/utils/password.go << 'EOF'
package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword gera um hash da senha
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifica se a senha corresponde ao hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
EOF

# Criar arquivo pagination.go
cat > internal/utils/pagination.go << 'EOF'
package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Pagination representa os parâmetros de paginação
type Pagination struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Sort      string `json:"sort"`
	Order     string `json:"order"`
	TotalRows int64  `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

// GetPaginationParams extrai os parâmetros de paginação da requisição
func GetPaginationParams(c *gin.Context) Pagination {
	// Valores padrão
	page := 1
	limit := 10
	sort := "created_at"
	order := "desc"

	// Extrair parâmetros da query
	if pageParam := c.Query("page"); pageParam != "" {
		pageInt, err := strconv.Atoi(pageParam)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		limitInt, err := strconv.Atoi(limitParam)
		if err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		sort = sortParam
	}

	if orderParam := c.Query("order"); orderParam != "" {
		if orderParam == "asc" || orderParam == "desc" {
			order = orderParam
		}
	}

	return Pagination{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}

// Paginate aplica a paginação a uma consulta GORM
func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) (*gorm.DB, error) {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	pagination.TotalPages = int(totalRows) / pagination.Limit
	if int(totalRows)%pagination.Limit > 0 {
		pagination.TotalPages++
	}

	offset := (pagination.Page - 1) * pagination.Limit
	return db.Offset(offset).Limit(pagination.Limit).Order(pagination.Sort + " " + pagination.Order), nil
}
EOF

# Criar arquivo response.go
cat > internal/utils/response.go << 'EOF'
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response representa a estrutura padrão de resposta da API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// SuccessResponse envia uma resposta de sucesso
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse envia uma resposta de erro
func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// ValidationErrorResponse envia uma resposta de erro de validação
func ValidationErrorResponse(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Data:    errors,
	})
}
EOF

# Criar arquivo auth.go (middleware)
cat > internal/api/middlewares/auth.go << 'EOF'
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/simple-erp-service/config"
	"github.com/yourusername/simple-erp-service/internal/utils"
)

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Token não fornecido")
			c.Abort()
			return
		}

		// Extrair token do cabeçalho
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Formato de token inválido")
			c.Abort()
			return
		}

		// Validar token
		claims, err := utils.ValidateToken(parts[1], cfg)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", err.Error())
			c.Abort()
			return  http.StatusUnauthorized, "Não autorizado", err.Error())
			c.Abort()
			return
		}

		// Armazenar claims no contexto
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}
EOF

# Criar arquivo permission.go (middleware)
cat > internal/api/middlewares/permission.go << 'EOF'
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/simple-erp-service/internal/utils"
)

// RequirePermission verifica se o usuário tem a permissão necessária
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se o usuário está autenticado
		permissions, exists := c.Get("permissions")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Usuário não autenticado")
			c.Abort()
			return
		}

		// Verificar se o usuário tem a permissão necessária
		userPermissions, ok := permissions.([]string)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro interno", "Erro ao verificar permissões")
			c.Abort()
			return
		}

		// Verificar se o usuário é admin (tem todas as permissões)
		role, exists := c.Get("role")
		if exists && role.(string) == "ADMIN" {
			c.Next()
			return
		}

		// Verificar se o usuário tem a permissão específica
		hasPermission := false
		for _, p := range userPermissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			utils.ErrorResponse(c, http.StatusForbidden, "Acesso negado", "Permissão insuficiente")
			c.Abort()
			return
		}

		c.Next()
	}
}
EOF

# Criar arquivo logger.go (middleware)
cat > internal/api/middlewares/logger.go << 'EOF'
package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/simple-erp-service/internal/models"
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
				UserID:    userID.(uint),
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
EOF

# Criar arquivo auth_service.go
cat > internal/service/auth_service.go << 'EOF'
package service

import (
	"errors"
	"time"

	"github.com/yourusername/simple-erp-service/config"
	"github.com/yourusername/simple-erp-service/internal/models"
	"github.com/yourusername/simple-erp-service/internal/utils"
	"gorm.io/gorm"
)

// AuthService gerencia a autenticação de usuários
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService cria um novo serviço de autenticação
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	User         models.UserResponse `json:"user"`
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	ExpiresIn    int                 `json:"expires_in"`
}

// Login autentica um usuário e retorna tokens JWT
func (s *AuthService) Login(username, password string) (*LoginResponse, error) {
	var user models.User

	// Buscar usuário pelo username
	result := s.db.Preload("Role.Permissions").Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, result.Error
	}

	// Verificar se o usuário está ativo
	if !user.IsActive {
		return nil, errors.New("usuário inativo")
	}

	// Verificar senha
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("senha incorreta")
	}

	// Extrair permissões
	var permissions []string
	if user.Role != nil {
		for _, perm := range user.Role.Permissions {
			permissions = append(permissions, perm.Name)
		}
	}

	// Gerar tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.RoleID, user.Role.Name, permissions, s.cfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, s.cfg)
	if err != nil {
		return nil, err
	}

	// Atualizar último login
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	return &LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.cfg.JWT.AccessTokenExp.Minutes()),
	}, nil
}

// RefreshToken renova o token de acesso usando um token de refresh
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// Validar token de refresh
	claims, err := utils.ValidateToken(refreshToken, s.cfg)
	if err != nil {
		return nil, err
	}

	// Buscar usuário
	var user models.User
	result := s.db.Preload("Role.Permissions").Where("username = ?", claims.Subject).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Verificar se o usuário está ativo
	if !user.IsActive {
		return nil, errors.New("usuário inativo")
	}

	// Extrair permissões
	var permissions []string
	if user.Role != nil {
		for _, perm := range user.Role.Permissions {
			permissions = append(permissions, perm.Name)
		}
	}

	// Gerar novo token de acesso
	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.RoleID, user.Role.Name, permissions, s.cfg)
	if err != nil {
		return nil, err
	}

	// Gerar novo token de refresh
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, s.cfg)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         user.ToResponse(),
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.cfg.JWT.AccessTokenExp.Minutes()),
	}, nil
}

// GetUserByID busca um usuário pelo ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := s.db.Preload("Role.Permissions").First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
EOF

# Criar arquivo user_service.go
cat > internal/service/user_service.go << 'EOF'
package service

import (
	"errors"

	"github.com/yourusername/simple-erp-service/internal/models"
	"github.com/yourusername/simple-erp-service/internal/utils"
	"gorm.io/gorm"
)

// UserService gerencia operações relacionadas a usuários
type UserService struct {
	db *gorm.DB
}

// NewUserService cria um novo serviço de usuários
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetUsers retorna uma lista paginada de usuários
func (s *UserService) GetUsers(pagination *utils.Pagination) ([]models.User, error) {
	var users []models.User

	query := s.db.Model(&models.User{}).Preload("Role")
	query, err := utils.Paginate(&models.User{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID busca um usuário pelo ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser cria um novo usuário
func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// Verificar se o username já existe
	var count int64
	s.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("username já está em uso")
	}

	// Verificar se o email já existe
	s.db.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("email já está em uso")
	}

	// Verificar se o perfil existe
	var role models.Role
	if err := s.db.First(&role, req.RoleID).Error; err != nil {
		return nil, errors.New("perfil não encontrado")
	}

	// Gerar hash da senha
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Criar usuário
	user := models.User{
		Username:     req.Username,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Carregar o perfil
	s.db.Preload("Role").First(&user, user.ID)

	return &user, nil
}

// UpdateUser atualiza um usuário existente
func (s *UserService) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	// Buscar usuário
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Verificar se o email já está em uso por outro usuário
	if req.Email != "" && req.Email != user.Email {
		var count int64
		s.db.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, id).Count(&count)
		if count > 0 {
			return nil, errors.New("email já está em uso")
		}
	}

	// Verificar se o perfil existe
	if req.RoleID != 0 && req.RoleID != user.RoleID {
		var role models.Role
		if err := s.db.First(&role, req.RoleID).Error; err != nil {
			return nil, errors.New("perfil não encontrado")
		}
	}

	// Atualizar campos
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.RoleID != 0 {
		user.RoleID = req.RoleID
	}

	// Salvar alterações
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	// Carregar o perfil
	s.db.Preload("Role").First(&user, user.ID)

	return &user, nil
}

// DeleteUser desativa um usuário
func (s *UserService) DeleteUser(id uint) error {
	// Buscar usuário
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	// Desativar usuário
	user.IsActive = false
	return s.db.Save(&user).Error
}

// ChangePassword altera a senha de um usuário
func (s *UserService) ChangePassword(id uint, currentPassword, newPassword string) error {
	// Buscar usuário
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	// Verificar senha atual
	if !utils.CheckPasswordHash(currentPassword, user.PasswordHash) {
		return errors.New("senha atual incorreta")
	}

	// Gerar hash da nova senha
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Atualizar senha
	user.PasswordHash = passwordHash
	return s.db.Save(&user).Error
}
EOF

# Criar arquivo auth.go (handler)
cat > internal/api/handlers/auth.go << 'EOF'
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/simple-erp-service/config"
	"github.com/yourusername/simple-erp-service/internal/service"
	"github.com/yourusername/simple-erp-service/internal/utils"
	"gorm.io/gorm"
)

// AuthHandler gerencia as requisições de autenticação
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler cria um novo handler de autenticação
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(db, cfg),
	}
}

// LoginRequest representa os dados de login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest representa os dados para renovar o token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login autentica um usuário
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados de login inválidos", err.Error())
		return
	}

	response, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Falha na autenticação", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login realizado com sucesso", response, nil)
}

// RefreshToken renova o token de acesso
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Token de refresh inválido", err.Error())
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Falha ao renovar token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token renovado com sucesso", response, nil)
}

// Logout realiza o logout do usuário
func (h *AuthHandler) Logout(c *gin.Context) {
	// Nota: Como estamos usando JWT, o logout é gerenciado pelo cliente
	// O servidor não precisa fazer nada além de retornar uma resposta de sucesso
	utils.SuccessResponse(c, http.StatusOK, "Logout realizado com sucesso", nil, nil)
}

// GetMe retorna informações do usuário logado
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Não autorizado", "Usuário não autenticado")
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar usuário", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário encontrado", user.ToResponse(), nil)
}
EOF

# Criar arquivo users.go (handler)
cat > internal/api/handlers/users.go << 'EOF'
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/simple-erp-service/internal/models"
	"github.com/yourusername/simple-erp-service/internal/service"
	"github.com/yourusername/simple-erp-service/internal/utils"
	"gorm.io/gorm"
)

// UserHandler gerencia as requisições relacionadas a usuários
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler cria um novo handler de usuários
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(db),
	}
}

// GetUsers retorna uma lista paginada de usuários
func (h *UserHandler) GetUsers(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	users, err := h.userService.GetUsers(&pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar usuários", err.Error())
		return
	}

	// Converter para resposta
	var response []models.UserResponse
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuários encontrados", response, pagination)
}

// GetUser retorna um usuário específico
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário encontrado", user.ToResponse(), nil)
}

// CreateUser cria um novo usuário
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar usuário", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Usuário criado com sucesso", user.ToResponse(), nil)
}

// UpdateUser atualiza um usuário existente
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar usuário", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário atualizado com sucesso", user.ToResponse(), nil)
}

// DeleteUser desativa um usuário
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao desativar usuário", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário desativado com sucesso", nil, nil)
}

// ChangePassword altera a senha de um usuário
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	// Verificar se o usuário está alterando sua própria senha
	userID, exists := c.Get("userID")
	if !exists || userID.(uint) != uint(id) {
		// Verificar se o usuário tem permissão para alterar senhas de outros usuários
		role, exists := c.Get("role")
		if !exists || role.(string) != "ADMIN" {
			utils.ErrorResponse(c, http.StatusForbidden, "Acesso negado", "Você não tem permissão para alterar a senha de outro usuário")
			return
		}
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	if err := h.userService.ChangePassword(uint(id), req.CurrentPassword, req.NewPassword); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao alterar senha", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Senha alterada com sucesso", nil, nil)
}
EOF

# Criar arquivo auth.go (routes)
cat > internal/api/routes/auth.go << 'EOF'
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/simple-erp-service/config"
	"github.com/yourusername/simple-erp-service/internal/api/handlers"
	"github.com/yourusername/simple-erp-service/internal/api/middlewares"
	"gorm.io/gorm"
)

// SetupAuthRoutes configura as rotas de autenticação
func SetupAuthRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
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
EOF

# Criar arquivo users.go (routes)
cat > internal/api/routes/users.go << 'EOF'
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/simple-erp-service/config"
	"github.com/yourusername/simple-erp-service/internal/api/handlers"
	"github.com/yourusername/simple-erp-service/internal/api/middlewares"
	"gorm.io/gorm"
)

// SetupUserRoutes configura as rotas de usuários
func SetupUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userHandler := handlers.NewUserHandler(db)

	// Obter configuração para middleware de autenticação
	cfg, _ := config.Load()

	// Grupo de rotas de usuários (todas protegidas)
	users := router.Group("/users")
	users.Use(middlewares.AuthMiddleware(cfg))
	{
		users.GET("", middlewares.RequirePermission("users.view"), userHandler.GetUsers)
		users.GET("/:id", middlewares.RequirePermission("users.view"), userHandler.GetUser)
		users.POST("", middlewares.RequirePermission("users.create"), userHandler.CreateUser)
		users.PUT("/:id", middlewares.RequirePermission("users.edit"), userHandler.UpdateUser)
		users.DELETE("/:id", middlewares.RequirePermission("users.delete"), userHandler.DeleteUser)
		users.PUT("/:id/password", userHandler.ChangePassword) // Permissão verificada no handler
	}
}
EOF

# Criar arquivos de placeholder para as outras rotas
for route in roles products inventory customers suppliers sales purchases financial dashboard system; do
    cat > internal/api/routes/${route}.go << EOF
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup${route^}Routes configura as rotas de ${route}
func Setup${route^}Routes(router *gin.RouterGroup, db *gorm.DB) {
	// TODO: Implementar rotas de ${route}
}
EOF
done

# Criar arquivo README.md
cat > README.md << 'EOF'
# Simple ERP Service

Backend em Golang para um sistema ERP simples, utilizando Gin e GORM.

## Requisitos

- Go 1.21 ou superior
- PostgreSQL 12 ou superior

## Configuração

1. Clone o repositório
2. Configure o arquivo `.env` com suas credenciais de banco de dados
3. Execute o script SQL para criar o banco de dados e tabelas (disponível em `migrations/create_database.sql`)
4. Execute o comando `go mod tidy` para instalar as dependências
5. Execute o comando `go run cmd/api/main.go` para iniciar o servidor

## Estrutura do Projeto
