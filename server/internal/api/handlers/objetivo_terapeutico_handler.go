package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msd-service/server/internal/models"
	"msd-service/server/internal/service"
)

// ObjetivoTerapeuticoHandler gerencia as requisições HTTP relacionadas a objetivos terapêuticos
type ObjetivoTerapeuticoHandler struct {
	service *service.ObjetivoTerapeuticoService
}

// NewObjetivoTerapeuticoHandler cria uma nova instância de ObjetivoTerapeuticoHandler
func NewObjetivoTerapeuticoHandler(service *service.ObjetivoTerapeuticoService) *ObjetivoTerapeuticoHandler {
	return &ObjetivoTerapeuticoHandler{service: service}
}

// CreateObjetivo godoc
// @Summary Criar um novo objetivo terapêutico
// @Description Cria um novo objetivo terapêutico com os dados fornecidos
// @Tags objetivos
// @Accept json
// @Produce json
// @Param objetivo body models.ObjetivoTerapeutico true "Dados do objetivo terapêutico"
// @Success 201 {object} models.ObjetivoTerapeutico
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/objetivos [post]
func (h *ObjetivoTerapeuticoHandler) CreateObjetivo(c *gin.Context) {
	var objetivo models.ObjetivoTerapeutico
	if err := c.ShouldBindJSON(&objetivo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateObjetivo(c.Request.Context(), &objetivo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetObjetivo godoc
// @Summary Obter um objetivo terapêutico pelo ID
// @Description Retorna os detalhes de um objetivo terapêutico específico
// @Tags objetivos
// @Accept json
// @Produce json
// @Param id path string true "ID do objetivo terapêutico"
// @Success 200 {object} models.ObjetivoTerapeutico
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Objetivo terapêutico não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/objetivos/{id} [get]
func (h *ObjetivoTerapeuticoHandler) GetObjetivo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	objetivo, err := h.service.GetObjetivo(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrObjetivoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Objetivo terapêutico não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, objetivo)
}

// UpdateObjetivo godoc
// @Summary Atualizar um objetivo terapêutico
// @Description Atualiza os dados de um objetivo terapêutico existente
// @Tags objetivos
// @Accept json
// @Produce json
// @Param id path string true "ID do objetivo terapêutico"
// @Param objetivo body models.ObjetivoTerapeutico true "Dados do objetivo terapêutico"
// @Success 200 {object} models.ObjetivoTerapeutico
// @Failure 400 {object} map[string]string "ID inválido ou erro de validação"
// @Failure 404 {object} map[string]string "Objetivo terapêutico não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/objetivos/{id} [put]
func (h *ObjetivoTerapeuticoHandler) UpdateObjetivo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var objetivo models.ObjetivoTerapeutico
	if err := c.ShouldBindJSON(&objetivo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objetivo.ID = id

	result, err := h.service.UpdateObjetivo(c.Request.Context(), &objetivo)
	if err != nil {
		if err == service.ErrObjetivoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Objetivo terapêutico não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteObjetivo godoc
// @Summary Excluir um objetivo terapêutico
// @Description Exclui um objetivo terapêutico pelo ID (soft delete)
// @Tags objetivos
// @Accept json
// @Produce json
// @Param id path string true "ID do objetivo terapêutico"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Objetivo terapêutico não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/objetivos/{id} [delete]
func (h *ObjetivoTerapeuticoHandler) DeleteObjetivo(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.service.DeleteObjetivo(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrObjetivoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Objetivo terapêutico não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListObjetivos godoc
// @Summary Listar objetivos terapêuticos
// @Description Retorna uma lista paginada de objetivos terapêuticos
// @Tags objetivos
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de objetivos terapêuticos e metadados de paginação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/objetivos [get]
func (h *ObjetivoTerapeuticoHandler) ListObjetivos(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page &lt; 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize &lt; 1 || pageSize > 100 {
		pageSize = 10
	}

	objetivos, total, err := h.service.ListObjetivos(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       objetivos,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// ListObjetivosByPaciente godoc
// @Summary Listar objetivos terapêuticos de um paciente
// @Description Retorna uma lista paginada de objetivos terapêuticos de um paciente específico
// @Tags objetivos
// @Accept json
// @Produce json
// @Param paciente_id path string true "ID do paciente"
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de objetivos terapêuticos e metadados de paginação"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes/{paciente_id}/objetivos [get]
func (h *ObjetivoTerapeuticoHandler) ListObjetivosByPaciente(c *gin.Context) {
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

	objetivos, total, err := h.service.ListObjetivosByPaciente(c.Request.Context(), pacienteID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       objetivos,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
