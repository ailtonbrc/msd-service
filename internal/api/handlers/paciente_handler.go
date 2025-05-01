package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"clinica_server/internal/auth"
	"clinica_server/internal/models"
	"clinica_server/internal/service"
)

// PacienteHandler define os handlers para as rotas de pacientes
type PacienteHandler struct {
	pacienteService service.PacienteService
	logger          *zap.Logger
}

// NewPacienteHandler cria uma nova instância de PacienteHandler
func NewPacienteHandler(pacienteService service.PacienteService, logger *zap.Logger) *PacienteHandler {
	return &PacienteHandler{
		pacienteService: pacienteService,
		logger:          logger,
	}
}

// Create godoc
// @Summary Criar um novo paciente
// @Description Cria um novo paciente no sistema
// @Tags pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param paciente body models.CreatePacienteRequest true "Dados do paciente"
// @Success 201 {object} models.PacienteResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pacientes [post]
func (h *PacienteHandler) Create(c *gin.Context) {
	// Obter contexto com informações do usuário
	ctx := c.Request.Context()

	// Validar e decodificar o corpo da requisição
	var req models.CreatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	// Criar paciente
	paciente, err := h.pacienteService.Create(ctx, &req)
	if err != nil {
		h.logger.Error("Erro ao criar paciente", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obter informações do usuário para logging
	userClaims, _ := auth.GetUserFromContext(ctx)
	h.logger.Info("Paciente criado com sucesso", 
		zap.Uint("id", paciente.ID), 
		zap.String("nome", paciente.Nome),
		zap.Uint("usuario_id", userClaims.UserID))
		
	c.JSON(http.StatusCreated, paciente)
}

// GetByID godoc
// @Summary Buscar um paciente pelo ID
// @Description Retorna os dados de um paciente específico
// @Tags pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param id path int true "ID do paciente"
// @Success 200 {object} models.PacienteResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pacientes/{id} [get]
func (h *PacienteHandler) GetByID(c *gin.Context) {
	// Obter contexto com informações do usuário
	ctx := c.Request.Context()

	// Obter ID do paciente
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar paciente
	paciente, err := h.pacienteService.GetByID(ctx, uint(id))
	if err != nil {
		h.logger.Error("Erro ao buscar paciente", zap.Error(err), zap.Uint64("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paciente)
}

// GetAll godoc
// @Summary Listar todos os pacientes
// @Description Retorna uma lista paginada de pacientes
// @Tags pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Param search query string false "Termo de busca"
// @Success 200 {object} models.PaginatedResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pacientes [get]
func (h *PacienteHandler) GetAll(c *gin.Context) {
	// Obter contexto com informações do usuário
	ctx := c.Request.Context()

	// Obter parâmetros de paginação
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")

	// Buscar pacientes
	result, err := h.pacienteService.GetAll(ctx, page, pageSize, search)
	if err != nil {
		h.logger.Error("Erro ao buscar pacientes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update godoc
// @Summary Atualizar um paciente
// @Description Atualiza os dados de um paciente existente
// @Tags pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param id path int true "ID do paciente"
// @Param paciente body models.UpdatePacienteRequest true "Dados do paciente"
// @Success 200 {object} models.PacienteResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pacientes/{id} [put]
func (h *PacienteHandler) Update(c *gin.Context) {
	// Obter contexto com informações do usuário
	ctx := c.Request.Context()

	// Obter ID do paciente
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Validar e decodificar o corpo da requisição
	var req models.UpdatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Erro ao decodificar requisição", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	// Atualizar paciente
	paciente, err := h.pacienteService.Update(ctx, uint(id), &req)
	if err != nil {
		h.logger.Error("Erro ao atualizar paciente", zap.Error(err), zap.Uint64("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obter informações do usuário para logging
	userClaims, _ := auth.GetUserFromContext(ctx)
	h.logger.Info("Paciente atualizado com sucesso", 
		zap.Uint("id", paciente.ID), 
		zap.String("nome", paciente.Nome),
		zap.Uint("usuario_id", userClaims.UserID))
		
	c.JSON(http.StatusOK, paciente)
}

// Delete godoc
// @Summary Excluir um paciente
// @Description Exclui um paciente do sistema
// @Tags pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param id path int true "ID do paciente"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pacientes/{id} [delete]
func (h *PacienteHandler) Delete(c *gin.Context) {
	// Obter contexto com informações do usuário
	ctx := c.Request.Context()

	// Obter ID do paciente
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Excluir paciente
	err = h.pacienteService.Delete(ctx, uint(id))
	if err != nil {
		h.logger.Error("Erro ao excluir paciente", zap.Error(err), zap.Uint64("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obter informações do usuário para logging
	userClaims, _ := auth.GetUserFromContext(ctx)
	h.logger.Info("Paciente excluído com sucesso", 
		zap.Uint64("id", id),
		zap.Uint("usuario_id", userClaims.UserID))
		
	c.JSON(http.StatusOK, gin.H{"message": "Paciente excluído com sucesso"})
}

// Search godoc
// @Summary Buscar pacientes com filtros avançados
// @Description Busca pacientes com filtros avançados
// @Tags pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param nome query string false "Nome do paciente"
// @Param cpf query string false "CPF do paciente"
// @Param telefone query string false "Telefone do paciente"
// @Param diagnostico query string false "Diagnóstico do paciente"
// @Param page query int false "Número da página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/pacientes/busca [get]
func (h *PacienteHandler) Search(c *gin.Context) {
	// Obter contexto com informações do usuário
	ctx := c.Request.Context()

	// Obter parâmetros de busca
	nome := c.Query("nome")
	cpf := c.Query("cpf")
	telefone := c.Query("telefone")
	diagnostico := c.Query("diagnostico")
	
	// Obter parâmetros de paginação
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Criar filtro de busca
	filtro := &models.PacienteFiltro{
		Nome:        nome,
		CPF:         cpf,
		Telefone:    telefone,
		Diagnostico: diagnostico,
		Page:        page,
		PageSize:    pageSize,
	}

	// Buscar pacientes
	result, err := h.pacienteService.Search(ctx, filtro)
	if err != nil {
		h.logger.Error("Erro ao buscar pacientes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}