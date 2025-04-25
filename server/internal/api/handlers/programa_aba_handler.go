package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msd-service/server/internal/models"
	"msd-service/server/internal/service"
)

// ProgramaABAHandler gerencia as requisições HTTP relacionadas a programas ABA
type ProgramaABAHandler struct {
	service *service.ProgramaABAService
}

// NewProgramaABAHandler cria uma nova instância de ProgramaABAHandler
func NewProgramaABAHandler(service *service.ProgramaABAService) *ProgramaABAHandler {
	return &ProgramaABAHandler{service: service}
}

// CreatePrograma godoc
// @Summary Criar um novo programa ABA
// @Description Cria um novo programa ABA com os dados fornecidos
// @Tags programas
// @Accept json
// @Produce json
// @Param programa body models.ProgramaABA true "Dados do programa ABA"
// @Success 201 {object} models.ProgramaABA
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/programas [post]
func (h *ProgramaABAHandler) CreatePrograma(c *gin.Context) {
	var programa models.ProgramaABA
	if err := c.ShouldBindJSON(&programa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreatePrograma(c.Request.Context(), &programa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetPrograma godoc
// @Summary Obter um programa ABA pelo ID
// @Description Retorna os detalhes de um programa ABA específico
// @Tags programas
// @Accept json
// @Produce json
// @Param id path string true "ID do programa ABA"
// @Success 200 {object} models.ProgramaABA
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Programa ABA não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/programas/{id} [get]
func (h *ProgramaABAHandler) GetPrograma(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	programa, err := h.service.GetPrograma(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrProgramaNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Programa ABA não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, programa)
}

// UpdatePrograma godoc
// @Summary Atualizar um programa ABA
// @Description Atualiza os dados de um programa ABA existente
// @Tags programas
// @Accept json
// @Produce json
// @Param id path string true "ID do programa ABA"
// @Param programa body models.ProgramaABA true "Dados do programa ABA"
// @Success 200 {object} models.ProgramaABA
// @Failure 400 {object} map[string]string "ID inválido ou erro de validação"
// @Failure 404 {object} map[string]string "Programa ABA não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/programas/{id} [put]
func (h *ProgramaABAHandler) UpdatePrograma(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var programa models.ProgramaABA
	if err := c.ShouldBindJSON(&programa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	programa.ID = id

	result, err := h.service.UpdatePrograma(c.Request.Context(), &programa)
	if err != nil {
		if err == service.ErrProgramaNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Programa ABA não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeletePrograma godoc
// @Summary Excluir um programa ABA
// @Description Exclui um programa ABA pelo ID (soft delete)
// @Tags programas
// @Accept json
// @Produce json
// @Param id path string true "ID do programa ABA"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Programa ABA não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/programas/{id} [delete]
func (h *ProgramaABAHandler) DeletePrograma(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.service.DeletePrograma(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrProgramaNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Programa ABA não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListProgramas godoc
// @Summary Listar programas ABA
// @Description Retorna uma lista paginada de programas ABA
// @Tags programas
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de programas ABA e metadados de paginação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/programas [get]
func (h *ProgramaABAHandler) ListProgramas(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page &lt; 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize &lt; 1 || pageSize > 100 {
		pageSize = 10
	}

	programas, total, err := h.service.ListProgramas(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       programas,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// ListProgramasByPaciente godoc
// @Summary Listar programas ABA de um paciente
// @Description Retorna uma lista paginada de programas ABA de um paciente específico
// @Tags programas
// @Accept json
// @Produce json
// @Param paciente_id path string true "ID do paciente"
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de programas ABA e metadados de paginação"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes/{paciente_id}/programas [get]
func (h *ProgramaABAHandler) ListProgramasByPaciente(c *gin.Context) {
	pacienteID, err := uuid.Parse(c.Param("paciente_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do paciente inválido"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page &lt; 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize &lt; 1 || pageSize > 100 {
		pageSize = 10
	}

	programas, total, err := h.service.ListProgramasByPaciente(c.Request.Context(), pacienteID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       programas,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
