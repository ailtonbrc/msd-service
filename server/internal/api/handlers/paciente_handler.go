package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msd-service/server/internal/models"
	"msd-service/server/internal/service"
)

// PacienteHandler gerencia as requisições HTTP relacionadas a pacientes
type PacienteHandler struct {
	service *service.PacienteService
}

// NewPacienteHandler cria uma nova instância de PacienteHandler
func NewPacienteHandler(service *service.PacienteService) *PacienteHandler {
	return &PacienteHandler{service: service}
}

// CreatePaciente godoc
// @Summary Criar um novo paciente
// @Description Cria um novo paciente com os dados fornecidos
// @Tags pacientes
// @Accept json
// @Produce json
// @Param paciente body models.CreatePacienteRequest true "Dados do paciente"
// @Success 201 {object} models.PacienteResponse
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes [post]
func (h *PacienteHandler) CreatePaciente(c *gin.Context) {
	var req models.CreatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreatePaciente(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetPaciente godoc
// @Summary Obter um paciente pelo ID
// @Description Retorna os detalhes de um paciente específico
// @Tags pacientes
// @Accept json
// @Produce json
// @Param id path string true "ID do paciente"
// @Success 200 {object} models.PacienteResponse
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Paciente não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes/{id} [get]
func (h *PacienteHandler) GetPaciente(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	resp, err := h.service.GetPaciente(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrPacienteNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Paciente não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdatePaciente godoc
// @Summary Atualizar um paciente
// @Description Atualiza os dados de um paciente existente
// @Tags pacientes
// @Accept json
// @Produce json
// @Param id path string true "ID do paciente"
// @Param paciente body models.UpdatePacienteRequest true "Dados do paciente"
// @Success 200 {object} models.PacienteResponse
// @Failure 400 {object} map[string]string "ID inválido ou erro de validação"
// @Failure 404 {object} map[string]string "Paciente não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes/{id} [put]
func (h *PacienteHandler) UpdatePaciente(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req models.UpdatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdatePaciente(c.Request.Context(), id, &req)
	if err != nil {
		if err == service.ErrPacienteNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Paciente não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeletePaciente godoc
// @Summary Excluir um paciente
// @Description Exclui um paciente pelo ID (soft delete)
// @Tags pacientes
// @Accept json
// @Produce json
// @Param id path string true "ID do paciente"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Paciente não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes/{id} [delete]
func (h *PacienteHandler) DeletePaciente(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.service.DeletePaciente(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrPacienteNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Paciente não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListPacientes godoc
// @Summary Listar pacientes
// @Description Retorna uma lista paginada de pacientes
// @Tags pacientes
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de pacientes e metadados de paginação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes [get]
func (h *PacienteHandler) ListPacientes(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	pacientes, total, err := h.service.ListPacientes(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       pacientes,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
