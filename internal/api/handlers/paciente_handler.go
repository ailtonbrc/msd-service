// internal/api/handlers/paciente_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"clinica_server/internal/auth"
	"clinica_server/internal/models"
	"clinica_server/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PacienteHandler gerencia as requisições HTTP relacionadas a pacientes
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

// @Summary Criar paciente
// @Description Cria um novo paciente
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param paciente body models.CreatePacienteRequest true "Dados do paciente"
// @Success 201 {object} models.Paciente
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes [post]
func (h *PacienteHandler) Create(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	// Logar a ação
	h.logger.Info("Criando novo paciente",
		zap.Uint("user_id", userClaims.UserID),
		zap.String("username", userClaims.Username),
	)

	var request models.CreatePacienteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Erro de validação ao criar paciente",
			zap.Error(err),
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Erro de validação",
			Message: err.Error(),
		})
		return
	}

	paciente, err := h.pacienteService.Create(c.Request.Context(), request)
	if err != nil {
		h.logger.Error("Erro ao criar paciente",
			zap.Error(err),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, response := handleServiceError(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Paciente criado com sucesso",
		zap.Uint("paciente_id", paciente.ID),
		zap.Uint("user_id", userClaims.UserID),
	)

	c.JSON(http.StatusCreated, paciente)
}

// @Summary Obter paciente por ID
// @Description Retorna os detalhes de um paciente específico
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param id path int true "ID do paciente"
// @Success 200 {object} models.PacienteDetailDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes/{id} [get]
func (h *PacienteHandler) GetByID(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("ID inválido ao buscar paciente",
			zap.String("id", c.Param("id")),
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID inválido",
			Message: "O ID do paciente deve ser um número inteiro válido",
		})
		return
	}

	h.logger.Info("Buscando paciente por ID",
		zap.Uint64("paciente_id", id),
		zap.Uint("user_id", userClaims.UserID),
	)

	paciente, err := h.pacienteService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Erro ao buscar paciente por ID",
			zap.Error(err),
			zap.Uint64("paciente_id", id),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, response := handleServiceError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, paciente)
}

// @Summary Atualizar paciente
// @Description Atualiza os dados de um paciente existente
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param id path int true "ID do paciente"
// @Param paciente body models.UpdatePacienteRequest true "Dados do paciente"
// @Success 200 {object} models.Paciente
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes/{id} [put]
func (h *PacienteHandler) Update(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("ID inválido ao atualizar paciente",
			zap.String("id", c.Param("id")),
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID inválido",
			Message: "O ID do paciente deve ser um número inteiro válido",
		})
		return
	}

	var request models.UpdatePacienteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Erro de validação ao atualizar paciente",
			zap.Error(err),
			zap.Uint64("paciente_id", id),
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Erro de validação",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("Atualizando paciente",
		zap.Uint64("paciente_id", id),
		zap.Uint("user_id", userClaims.UserID),
	)

	paciente, err := h.pacienteService.Update(c.Request.Context(), uint(id), request)
	if err != nil {
		h.logger.Error("Erro ao atualizar paciente",
			zap.Error(err),
			zap.Uint64("paciente_id", id),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, response := handleServiceError(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Paciente atualizado com sucesso",
		zap.Uint64("paciente_id", id),
		zap.Uint("user_id", userClaims.UserID),
	)

	c.JSON(http.StatusOK, paciente)
}

// @Summary Excluir paciente
// @Description Remove um paciente do sistema
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param id path int true "ID do paciente"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes/{id} [delete]
func (h *PacienteHandler) Delete(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("ID inválido ao excluir paciente",
			zap.String("id", c.Param("id")),
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID inválido",
			Message: "O ID do paciente deve ser um número inteiro válido",
		})
		return
	}

	h.logger.Info("Excluindo paciente",
		zap.Uint64("paciente_id", id),
		zap.Uint("user_id", userClaims.UserID),
	)

	err = h.pacienteService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Erro ao excluir paciente",
			zap.Error(err),
			zap.Uint64("paciente_id", id),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, response := handleServiceError(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Paciente excluído com sucesso",
		zap.Uint64("paciente_id", id),
		zap.Uint("user_id", userClaims.UserID),
	)

	c.Status(http.StatusNoContent)
}

// @Summary Listar pacientes
// @Description Retorna uma lista paginada de pacientes
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Itens por página" default(10)
// @Param nome query string false "Filtrar por nome"
// @Param cpf query string false "Filtrar por CPF"
// @Param diagnostico query string false "Filtrar por diagnóstico"
// @Param genero query string false "Filtrar por gênero"
// @Param cidade query string false "Filtrar por cidade"
// @Param estado query string false "Filtrar por estado"
// @Success 200 {object} models.PacienteListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes [get]
func (h *PacienteHandler) List(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	// Parâmetros de paginação
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Filtros
	filters := make(map[string]interface{})
	if nome := c.Query("nome"); nome != "" {
		filters["nome"] = nome
	}
	if cpf := c.Query("cpf"); cpf != "" {
		filters["cpf"] = cpf
	}
	if diagnostico := c.Query("diagnostico"); diagnostico != "" {
		filters["diagnostico"] = diagnostico
	}
	if genero := c.Query("genero"); genero != "" {
		filters["genero"] = genero
	}
	if cidade := c.Query("cidade"); cidade != "" {
		filters["cidade"] = cidade
	}
	if estado := c.Query("estado"); estado != "" {
		filters["estado"] = estado
	}

	h.logger.Info("Listando pacientes",
		zap.Int("page", page),
		zap.Int("limit", limit),
		zap.Any("filters", filters),
		zap.Uint("user_id", userClaims.UserID),
	)

	response, err := h.pacienteService.List(c.Request.Context(), page, limit, filters)
	if err != nil {
		h.logger.Error("Erro ao listar pacientes",
			zap.Error(err),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, errResponse := handleServiceError(err)
		c.JSON(statusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Buscar pacientes
// @Description Busca pacientes por um termo de busca
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param q query string true "Termo de busca"
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Itens por página" default(10)
// @Success 200 {object} models.PacienteListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes/search [get]
func (h *PacienteHandler) Search(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	query := c.Query("q")
	if query == "" {
		h.logger.Warn("Termo de busca vazio",
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Parâmetro inválido",
			Message: "O termo de busca (q) é obrigatório",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	h.logger.Info("Buscando pacientes",
		zap.String("query", query),
		zap.Int("page", page),
		zap.Int("limit", limit),
		zap.Uint("user_id", userClaims.UserID),
	)

	response, err := h.pacienteService.Search(c.Request.Context(), query, page, limit)
	if err != nil {
		h.logger.Error("Erro ao buscar pacientes",
			zap.Error(err),
			zap.String("query", query),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, errResponse := handleServiceError(err)
		c.JSON(statusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Buscar paciente por CPF
// @Description Retorna os detalhes de um paciente pelo CPF
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param cpf path string true "CPF do paciente"
// @Success 200 {object} models.PacienteDetailDTO
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes/cpf/{cpf} [get]
func (h *PacienteHandler) GetByCPF(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	cpf := c.Param("cpf")
	if cpf == "" {
		h.logger.Warn("CPF vazio",
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Parâmetro inválido",
			Message: "O CPF é obrigatório",
		})
		return
	}

	h.logger.Info("Buscando paciente por CPF",
		zap.String("cpf", cpf),
		zap.Uint("user_id", userClaims.UserID),
	)

	paciente, err := h.pacienteService.GetByCPF(c.Request.Context(), cpf)
	if err != nil {
		h.logger.Error("Erro ao buscar paciente por CPF",
			zap.Error(err),
			zap.String("cpf", cpf),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, response := handleServiceError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, paciente)
}

// @Summary Calcular idade do paciente
// @Description Retorna a idade calculada de um paciente
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param id path int true "ID do paciente"
// @Success 200 {object} IdadeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes/{id}/idade [get]
func (h *PacienteHandler) CalcularIdade(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("ID inválido ao calcular idade",
			zap.String("id", c.Param("id")),
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID inválido",
			Message: "O ID do paciente deve ser um número inteiro válido",
		})
		return
	}

	h.logger.Info("Calculando idade do paciente",
		zap.Uint64("paciente_id", id),
		zap.Uint("user_id", userClaims.UserID),
	)

	idade, err := h.pacienteService.CalcularIdade(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Erro ao calcular idade do paciente",
			zap.Error(err),
			zap.Uint64("paciente_id", id),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, response := handleServiceError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, IdadeResponse{
		Idade: idade,
	})
}

// @Summary Atualizar diagnóstico
// @Description Atualiza apenas o diagnóstico de um paciente
// @Tags Pacientes
// @Accept json
// @Produce json
// @Param Authorization header string true "Token JWT"
// @Param id path int true "ID do paciente"
// @Param diagnostico body DiagnosticoRequest true "Novo diagnóstico"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security JWT
// @Router /api/pacientes/{id}/diagnostico [patch]
func (h *PacienteHandler) AtualizarDiagnostico(c *gin.Context) {
	// Extrair informações do usuário do contexto
	userClaims, exists := auth.GetUserFromContext(c.Request.Context())
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Não autorizado",
			Message: "Usuário não autenticado",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Warn("ID inválido ao atualizar diagnóstico",
			zap.String("id", c.Param("id")),
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "ID inválido",
			Message: "O ID do paciente deve ser um número inteiro válido",
		})
		return
	}

	var request DiagnosticoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Erro de validação ao atualizar diagnóstico",
			zap.Error(err),
			zap.Uint64("paciente_id", id),
			zap.Uint("user_id", userClaims.UserID),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Erro de validação",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info("Atualizando diagnóstico do paciente",
		zap.Uint64("paciente_id", id),
		zap.Uint("user_id", userClaims.UserID),
	)

	err = h.pacienteService.AtualizarDiagnostico(c.Request.Context(), uint(id), request.Diagnostico)
	if err != nil {
		h.logger.Error("Erro ao atualizar diagnóstico do paciente",
			zap.Error(err),
			zap.Uint64("paciente_id", id),
			zap.Uint("user_id", userClaims.UserID),
		)
		statusCode, response := handleServiceError(err)
		c.JSON(statusCode, response)
		return
	}

	h.logger.Info("Diagnóstico do paciente atualizado com sucesso",
		zap.Uint64("paciente_id", id),
		zap.Uint("user_id", userClaims.UserID),
	)

	c.Status(http.StatusNoContent)
}

// ErrorResponse representa uma resposta de erro da API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// IdadeResponse representa a resposta da API para o cálculo de idade
type IdadeResponse struct {
	Idade int `json:"idade"`
}

// DiagnosticoRequest representa a requisição para atualizar o diagnóstico
type DiagnosticoRequest struct {
	Diagnostico string `json:"diagnostico" binding:"required"`
}

// handleServiceError mapeia erros do serviço para códigos HTTP e respostas de erro
func handleServiceError(err error) (int, ErrorResponse) {
	switch {
	case service.ErrNotFound.Error() == err.Error():
		return http.StatusNotFound, ErrorResponse{
			Error:   "Não encontrado",
			Message: "Paciente não encontrado",
		}
	case service.ErrInvalidInput.Error() == err.Error():
		return http.StatusBadRequest, ErrorResponse{
			Error:   "Entrada inválida",
			Message: err.Error(),
		}
	case service.ErrUnauthorized.Error() == err.Error():
		return http.StatusForbidden, ErrorResponse{
			Error:   "Acesso proibido",
			Message: "Você não tem permissão para realizar esta operação",
		}
	case service.ErrDuplicateResource.Error() == err.Error():
		return http.StatusConflict, ErrorResponse{
			Error:   "Recurso duplicado",
			Message: err.Error(),
		}
	default:
		return http.StatusInternalServerError, ErrorResponse{
			Error:   "Erro interno",
			Message: "Ocorreu um erro interno no servidor",
		}
	}
}
