package handlers

import (
	"net/http"
	"strconv"

	"clinica_server/internal/models"
	"clinica_server/internal/repository"
	"clinica_server/internal/service"
	"clinica_server/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UsuarioHandler gerencia as requisições relacionadas a usuários
type UsuarioHandler struct {
	usuarioService *service.UsuarioService
}

// NewUsuarioHandler cria um novo handler de usuários
func NewUsuarioHandler(db *gorm.DB) *UsuarioHandler {
	usuarioRepo := repository.NewUsuarioRepository(db)

	return &UsuarioHandler{
		usuarioService: service.NewUsuarioService(usuarioRepo),
	}
}

func (h *UsuarioHandler) BuscaUsuarios(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	usuarios, err := h.usuarioService.BuscaUsuarios(&pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar usuários", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuários encontrados", usuarios, nil)
}

func (h *UsuarioHandler) BuscaUsuario(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	usuario, err := h.usuarioService.BuscaUsuarioPorID(uint(id))
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário encontrado", usuario, nil)
}

func (h *UsuarioHandler) CreateUsuario(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	usuario, err := h.usuarioService.CreateUsuario(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar usuário", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Usuário criado com sucesso", usuario, nil)
}

func (h *UsuarioHandler) UpdateUsuario(c *gin.Context) {
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

	usuario, err := h.usuarioService.UpdateUsuario(uint(id), req)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário atualizado com sucesso", usuario, nil)
}

func (h *UsuarioHandler) ChangePassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	// TODO Precisa Definir como q vai tratar o perfil do usuário

	// Verificar se o usuário é admin ou está alterando a própria senha
	usuarioID, _ := c.Get("usuarioID")
	//role, _ := c.Get("role")
	//isAdmin := role == "ADMIN"
	isSelf := usuarioID.(uint) == uint(id)

	// Se não for admin e não for o próprio usuário, negar acesso
	if /*!isAdmin &&*/ !isSelf {
		utils.ErrorResponse(c, http.StatusForbidden, "Acesso negado", "Você não tem permissão para alterar a senha de outro usuário")
		return
	}

	if /*isAdmin &&*/ !isSelf {
		// Se for admin alterando a senha de outro usuário, não precisa da senha atual
		err = h.usuarioService.TrocaSenha(uint(id), "", req.NewPassword, true)
	} else {
		// Se for o próprio usuário ou admin alterando a própria senha, precisa da senha atual
		err = h.usuarioService.TrocaSenha(uint(id), req.CurrentPassword, req.NewPassword, false)
	}

	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else if err == utils.ErrInvalidCredentials {
			utils.ErrorResponse(c, http.StatusBadRequest, "Senha atual incorreta", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao alterar senha", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Senha alterada com sucesso", nil, nil)
}

func (h *UsuarioHandler) DeleteUsuario(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	// Impedir que um usuário exclua a si mesmo
	usuarioID, _ := c.Get("usuarioID")
	if usuarioID.(uint) == uint(id) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Operação inválida", "Você não pode excluir seu próprio usuário")
		return
	}

	if err := h.usuarioService.DeleteUsuario(uint(id)); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Usuário não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir usuário", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuário excluído com sucesso", nil, nil)
}
