package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msd-service/server/internal/models"
	"msd-service/server/internal/service"
)

// TerapiaHandler gerencia as requisições HTTP relacionadas a terapias
type TerapiaHandler struct {
	service *service.TerapiaService
}

// NewTerapiaHandler cria uma nova instância de TerapiaHandler
func NewTerapiaHandler(service *service.TerapiaService) *TerapiaHandler {
	return &TerapiaHandler{service: service}
}

// CreateTerapia godoc
// @Summary Criar uma nova terapia
// @Description Cria uma nova terapia com os dados fornecidos
// @Tags terapias
// @Accept json
// @Produce json
// @Param terapia body models.Terapia true "Dados da terapia"
// @Success 201 {object} models.Terapia
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/terapias [post]
func (h *TerapiaHandler) CreateTerapia(c *gin.Context) {
	var terapia models.Terapia
	if err := c.ShouldBindJSON(&terapia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateTerapia(c.Request.Context(), &terapia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetTerapia godoc
// @Summary Obter uma terapia pelo ID
// @Description Retorna os detalhes de uma terapia específica
// @Tags terapias
// @Accept json
// @Produce json
// @Param id path string true "ID da terapia"
// @Success 200 {object} models.Terapia
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Terapia não encontrada"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/terapias/{id} [get]
func (h *TerapiaHandler) GetTerapia(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	terapia, err := h.service.GetTerapia(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrTerapiaNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Terapia não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, terapia)
}

// UpdateTerapia godoc
// @Summary Atualizar uma terapia
// @Description Atualiza os dados de uma terapia existente
// @Tags terapias
// @Accept json
// @Produce json
// @Param id path string true "ID da terapia"
// @Param terapia body models.Terapia true "Dados da terapia"
// @Success 200 {object} models.Terapia
// @Failure 400 {object} map[string]string "ID inválido ou erro de validação"
// @Failure 404 {object} map[string]string "Terapia não encontrada"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/terapias/{id} [put]
func (h *TerapiaHandler) UpdateTerapia(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var terapia models.Terapia
	if err := c.ShouldBindJSON(&terapia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	terapia.ID = id

	result, err := h.service.UpdateTerapia(c.Request.Context(), &terapia)
	if err != nil {
		if err == service.ErrTerapiaNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Terapia não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteTerapia godoc
//  err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteTerapia godoc
// @Summary Excluir uma terapia
// @Description Exclui uma terapia pelo ID (soft delete)
// @Tags terapias
// @Accept json
// @Produce json
// @Param id path string true "ID da terapia"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Terapia não encontrada"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/terapias/{id} [delete]
func (h *TerapiaHandler) DeleteTerapia(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.service.DeleteTerapia(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrTerapiaNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Terapia não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTerapias godoc
// @Summary Listar terapias
// @Description Retorna uma lista paginada de terapias
// @Tags terapias
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de terapias e metadados de paginação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/terapias [get]
func (h *TerapiaHandler) ListTerapias(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	terapias, total, err := h.service.ListTerapias(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       terapias,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

Let me continue with the remaining handlers and routes to complete the project:

```go file="msd-service/server/internal/api/handlers/sessao_handler.go"
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"msd-service/server/internal/models"
	"msd-service/server/internal/service"
)

// SessaoHandler gerencia as requisições HTTP relacionadas a sessões
type SessaoHandler struct {
	service *service.SessaoService
}

// NewSessaoHandler cria uma nova instância de SessaoHandler
func NewSessaoHandler(service *service.SessaoService) *SessaoHandler {
	return &SessaoHandler{service: service}
}

// CreateSessao godoc
// @Summary Criar uma nova sessão
// @Description Cria uma nova sessão com os dados fornecidos
// @Tags sessoes
// @Accept json
// @Produce json
// @Param sessao body models.Sessao true "Dados da sessão"
// @Success 201 {object} models.Sessao
// @Failure 400 {object} map[string]string "Erro de validação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/sessoes [post]
func (h *SessaoHandler) CreateSessao(c *gin.Context) {
	var sessao models.Sessao
	if err := c.ShouldBindJSON(&sessao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateSessao(c.Request.Context(), &sessao)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetSessao godoc
// @Summary Obter uma sessão pelo ID
// @Description Retorna os detalhes de uma sessão específica
// @Tags sessoes
// @Accept json
// @Produce json
// @Param id path string true "ID da sessão"
// @Success 200 {object} models.Sessao
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Sessão não encontrada"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/sessoes/{id} [get]
func (h *SessaoHandler) GetSessao(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	sessao, err := h.service.GetSessao(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrSessaoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sessão não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessao)
}

// UpdateSessao godoc
// @Summary Atualizar uma sessão
// @Description Atualiza os dados de uma sessão existente
// @Tags sessoes
// @Accept json
// @Produce json
// @Param id path string true "ID da sessão"
// @Param sessao body models.Sessao true "Dados da sessão"
// @Success 200 {object} models.Sessao
// @Failure 400 {object} map[string]string "ID inválido ou erro de validação"
// @Failure 404 {object} map[string]string "Sessão não encontrada"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/sessoes/{id} [put]
func (h *SessaoHandler) UpdateSessao(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var sessao models.Sessao
	if err := c.ShouldBindJSON(&sessao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sessao.ID = id

	result, err := h.service.UpdateSessao(c.Request.Context(), &sessao)
	if err != nil {
		if err == service.ErrSessaoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sessão não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteSessao godoc
// @Summary Excluir uma sessão
// @Description Exclui uma sessão pelo ID (soft delete)
// @Tags sessoes
// @Accept json
// @Produce json
// @Param id path string true "ID da sessão"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 404 {object} map[string]string "Sessão não encontrada"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/sessoes/{id} [delete]
func (h *SessaoHandler) DeleteSessao(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.service.DeleteSessao(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrSessaoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sessão não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSessoes godoc
// @Summary Listar sessões
// @Description Retorna uma lista paginada de sessões
// @Tags sessoes
// @Accept json
// @Produce json
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de sessões e metadados de paginação"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/sessoes [get]
func (h *SessaoHandler) ListSessoes(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page &lt; 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize &lt; 1 || pageSize > 100 {
		pageSize = 10
	}

	sessoes, total, err := h.service.ListSessoes(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       sessoes,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// ListSessoesByPaciente godoc
// @Summary Listar sessões de um paciente
// @Description Retorna uma lista paginada de sessões de um paciente específico
// @Tags sessoes
// @Accept json
// @Produce json
// @Param paciente_id path string true "ID do paciente"
// @Param page query int false "Número da página (padrão: 1)"
// @Param page_size query int false "Tamanho da página (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de sessões e metadados de paginação"
// @Failure 400 {object} map[string]string "ID inválido"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Security BearerAuth
// @Router /api/v1/pacientes/{paciente_id}/sessoes [get]
func (h *SessaoHandler) ListSessoesByPaciente(c *gin.Context) {
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

	sessoes, total, err := h.service.ListSessoesByPaciente(c.Request.Context(), pacienteID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       sessoes,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
