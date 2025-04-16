package handlers

import (
	"clinica_server/internal/models"
	"clinica_server/internal/service"
	"clinica_server/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PacienteHandler struct {
	pacienteService *service.PacienteService
}

func NewPacienteHandler(db *gorm.DB) *PacienteHandler {
	return &PacienteHandler{
		pacienteService: service.NewPacienteService(db),
	}
}

func (h *PacienteHandler) GetPacientes(c *gin.Context) {
	items, err := h.pacienteService.GetPacientes()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar pacientes", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Pacientes encontrados", items, nil)
}

func (h *PacienteHandler) GetPaciente(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}
	item, err := h.pacienteService.GetPacienteByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Paciente não encontrado", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Paciente encontrado", item, nil)
}

func (h *PacienteHandler) CreatePaciente(c *gin.Context) {
	var req models.CreatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}
	item, err := h.pacienteService.CreatePaciente(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar paciente", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Paciente criado com sucesso", item, nil)
}

func (h *PacienteHandler) UpdatePaciente(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}
	var req models.UpdatePacienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}
	item, err := h.pacienteService.UpdatePaciente(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar paciente", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Paciente atualizado com sucesso", item, nil)
}

func (h *PacienteHandler) DeletePaciente(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}
	if err := h.pacienteService.DeletePaciente(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir paciente", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Paciente excluído com sucesso", nil, nil)
}