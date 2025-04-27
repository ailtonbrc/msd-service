package handlers

import (
	"net/http"

	"clinica_server/internal/service"
	"clinica_server/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler gerencia as requisições de autenticação
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler cria um novo handler de autenticação
// func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
// 	return &AuthHandler{
// 		authService: service.NewAuthService(db, cfg),
// 	}
// }

// NewAuthHandler cria um novo handler de autenticação
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

// LoginRequest representa os dados de login
type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Senha string `json:"senha" binding:"required"`
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

	response, err := h.authService.Login(req.Email, req.Senha)
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

	utils.SuccessResponse(c, http.StatusOK, "Usuário encontrado", user.ToDTO(), nil)
}
