package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msd-service/server/internal/models"
	"msd-service/server/internal/service"
)

// ComportamentoAlvoHandler gerencia as requisições HTTP relacionadas a comportamentos alvo
type ComportamentoAlvoHandler struct {
	service *service.ComportamentoAlvoService
}

// NewComportamentoAlvoHandler cria uma nova instância de ComportamentoAlvoHandler
func NewComportamentoAlvoHandler(service *service.ComportamentoAlvoService) *ComportamentoAlvoHandler {
	return &ComportamentoAlvoHandler{service: service}
}

// CreateComportamento godoc
// @Summary Criar um novo comportamento alvo
// @Description Cria um novo comportamento alvo com os dados fornecidos
// @Tags comportamentos
// @Accept json
// @Produce json
// @Param comportamento body models.ComportamentoAlvo true "Dados do comportamento alvo"
// @Success 201 {object} models.ComportamentoAlvo
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/comportamentos [post]
func (h *ComportamentoAlvoHandler) CreateComportamento(c *gin.Context) {
	var comportamento models.ComportamentoAlvo
	if err := c.ShouldBindJSON(&comportamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateComportamento(c.Request.Context(), &comportamento)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetComportamento godoc
// @Summary Obter um comportamento alvo pelo ID
// @Description Retorna os detalhes de um comportamento alvo específico
// @Tags comportamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do comportamento alvo"
// @Success 200 {object} models.ComportamentoAlvo
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Comportamento alvo não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/comportamentos/{id} [get]
func (h *ComportamentoAlvoHandler) GetComportamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	comportamento, err := h.service.GetComportamento(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrComportamentoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comportamento alvo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comportamento)
}

// UpdateComportamento godoc
// @Summary Atualizar um comportamento alvo
// @Description Atualiza os dados de um comportamento alvo existente
// @Tags comportamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do comportamento alvo"
// @Param comportamento body models.ComportamentoAlvo true "Dados do comportamento alvo"
// @Success 200 {object} models.ComportamentoAlvo
// @Failure 400 {object} map[string]string "ID inválido ou erro de validação"
// @Failure 404 {object} map[string]string "Comportamento alvo não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/comportamentos/{id} [put]
func (h *ComportamentoAlvoHandler) UpdateComportamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var comportamento models.ComportamentoAlvo
	if err := c.ShouldBindJSON(&comportamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comportamento.ID = id

	result, err := h.service.UpdateComportamento(c.Request.Context(), &comportamento)
	if err != nil {
		if err == service.ErrComportamentoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comportamento alvo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteComportamento godoc
// @Summary Excluir um comportamento alvo
// @Description Exclui um comportamento alvo pelo ID (soft delete)
// @Tags comportamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do comportamento alvo"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Comportamento alvo não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/comportamentos/{id} [delete]
func (h *ComportamentoAlvoHandler) DeleteComportamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.service.DeleteComportamento(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrComportamentoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comportamento alvo não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListComportamentos godoc
// @Summary Listar comportamentos alvo
// @Description Retorna uma lista paginada de comportamentos alvo
// @Tags comportamentos
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de comportamentos alvo e metadados de paginação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/comportamentos [get]
func (h *ComportamentoAlvoHandler) ListComportamentos(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page &lt; 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize &lt; 1 || pageSize > 100 {
		pageSize = 10
	}

	comportamentos, total, err := h.service.ListComportamentos(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       comportamentos,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// ListComportamentosByPaciente godoc
// @Summary Listar comportamentos alvo de um paciente
// @Description Retorna uma lista paginada de comportamentos alvo de um paciente específico
// @Tags comportamentos
// @Accept json
// @Produce json
// @Param paciente_id path string true "ID do paciente"
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de comportamentos alvo e metadados de paginação"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes/{paciente_id}/comportamentos [get]
func (h *ComportamentoAlvoHandler) ListComportamentosByPaciente(c *gin.Context) {
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

	comportamentos, total, err := h.service.ListComportamentosByPaciente(c.Request.Context(), pacienteID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       comportamentos,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
