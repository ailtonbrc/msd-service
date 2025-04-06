package handlers

// import (
//     "net/http"
//     "strconv"

//     "clinica_server/internal/models"
//     "clinica_server/internal/service"
//     "clinica_server/internal/utils"

//     "github.com/gin-gonic/gin"
//     "gorm.io/gorm"
// )

// type ClinicaHandler struct {
//     clinicaService *service.ClinicaService
// }

// func NewClinicaHandler(db *gorm.DB) *ClinicaHandler {
//     return &ClinicaHandler{
//         clinicaService: service.NewClinicaService(db),
//     }
// }

// func (h *ClinicaHandler) GetClinicas(c *gin.Context) {
//     pagination := utils.GetPaginationParams(c)
//     items, err := h.clinicaService.GetClinicas(&pagination)
//     if err != nil {
//         utils.ErrorResponse(c, http.StatusInternalServerError, "Erro ao buscar clinicas", err.Error())
//         return
//     }
//     utils.SuccessResponse(c, http.StatusOK, "Clinicas encontrados", items, nil)
// }

// func (h *ClinicaHandler) GetClinica(c *gin.Context) {
//     id, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
//         return
//     }
//     item, err := h.clinicaService.GetClinicaByID(uint(id))
//     if err != nil {
//         utils.ErrorResponse(c, http.StatusNotFound, "Clinica não encontrado", err.Error())
//         return
//     }
//     utils.SuccessResponse(c, http.StatusOK, "Clinica encontrado", item, nil)
// }

// func (h *ClinicaHandler) CreateClinica(c *gin.Context) {
//     var req models.CreateClinicaRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
//         return
//     }
//     item, err := h.clinicaService.CreateClinica(req)
//     if err != nil {
//         utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao criar clinica", err.Error())
//         return
//     }
//     utils.SuccessResponse(c, http.StatusCreated, "Clinica criado com sucesso", item, nil)
// }

// func (h *ClinicaHandler) UpdateClinica(c *gin.Context) {
//     id, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
//         return
//     }
//     var req models.UpdateClinicaRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         utils.ValidationErrorResponse(c, "Dados inválidos", err.Error())
//         return
//     }
//     item, err := h.clinicaService.UpdateClinica(uint(id), req)
//     if err != nil {
//         utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao atualizar clinica", err.Error())
//         return
//     }
//     utils.SuccessResponse(c, http.StatusOK, "Clinica atualizado com sucesso", item, nil)
// }

// func (h *ClinicaHandler) DeleteClinica(c *gin.Context) {
//     id, err := strconv.ParseUint(c.Param("id"), 10, 32)
//     if err != nil {
//         utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido", err.Error())
//         return
//     }
//     if err := h.clinicaService.DeleteClinica(uint(id)); err != nil {
//         utils.ErrorResponse(c, http.StatusBadRequest, "Erro ao excluir clinica", err.Error())
//         return
//     }
//     utils.SuccessResponse(c, http.StatusOK, "Clinica excluído com sucesso", nil, nil)
// }
