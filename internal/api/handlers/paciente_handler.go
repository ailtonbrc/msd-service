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

// PacienteHandler gerencia as requisições relacionadas a paciente
type PacienteHandler struct {
	pacienteService *service.PacienteService
}

// NewPacienteHandler cria um novo handler de paciente
func NewPacienteHandler(db *gorm.DB) *PacienteHandler {
	pacienteRepo := repository.NewPacienteRepository(db)

	return &PacienteHandler{
		pacienteService: service.NewPacienteService(pacienteRepo),
	}
}

func (h *PacienteHandler) BuscaPacientes(c *gin.Context) {
	pagination := utils.GetPaginationParams(c)

	pacientes, err := h.pacienteService.BuscaPacientes(&pagination)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar pacientes", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pacientes encontrados", pacientes, nil)
}

func (h *PacienteHandler) BuscaPaciente(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	paciente, err := h.pacienteService.BuscaPacientePorID(uint(id))
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Paciente não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar paciente", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Paciente encontrado", paciente, nil)
}

func (h *PacienteHandler) CreatePaciente(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
		return
	}

	paciente, err := h.pacienteService.CreatePaciente(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar paciente", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Paciente criado com sucesso", paciente, nil)
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

	paciente, err := h.pacienteService.UpdatePaciente(uint(id), req)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Paciente não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar paciente", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Paciente atualizado com sucesso", paciente, nil)
}

func (h *PacienteHandler) DeletePaciente(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
		return
	}

	if err := h.pacienteService.DeletePaciente(uint(id)); err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Paciente não encontrado", err.Error())
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir paciente", err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Paciente excluído com sucesso", nil, nil)
}
