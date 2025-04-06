package models

import "clinica_server/internal/utils"

// PaginationDTO representa informações de paginação
type PaginationDTO struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Order      string `json:"order"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
}

// ToPaginationDTO converte um objeto de paginação para PaginationDTO
func ToPaginationDTO(pagination *utils.Pagination) *PaginationDTO {
	if pagination == nil {
		return nil
	}

	return &PaginationDTO{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Sort:       pagination.Sort,
		Order:      pagination.Order,
		TotalRows:  pagination.TotalRows,
		TotalPages: pagination.TotalPages,
	}
}
