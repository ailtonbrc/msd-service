package models

// UsuarioDTO representa os dados de um usuário que são seguros para enviar ao frontend
type UsuarioDTO struct {
	ID     uint   `json:"id"`
	Nome   string `json:"name"`
	Email  string `json:"email,omitempty"`
	Perfil string `json:"perfil,omitempty"`
	Ativo  bool   `json:"is_active"`
}

// UsuarioDetailDTO representa os dados detalhados de um usuário
type UsuarioDetalheDTO struct {
	ID                    uint   `json:"id"`
	Nome                  string `json:"nome"`
	Email                 string `json:"email,omitempty"`
	Perfil                string `json:"perfil,omitempty"`
	Ativo                 bool   `json:"ativo"`
	ClinicaID             uint   `json:"clinica_id,omitempty"`
	SupervisorID          uint   `json:"supervisor_id,omitempty"`
	DataInicioInatividade string `json:"data_inicio_inatividade,omitempty"`
	DataFimInatividade    string `json:"data_fim_inatividade,omitempty"`
	MotivoInatividade     string `json:"motivo_inatividade,omitempty"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
}

// UsuarioListDTO representa uma lista paginada de usuários
type UsuarioListDTO struct {
	Usuarios   []UsuarioDTO   `json:"usuarios"`
	Pagination *PaginationDTO `json:"pagination,omitempty"`
}

// ToDTO converte um modelo Usuario para UsuarioDTO
func (u *Usuario) ToDTO() UsuarioDTO {
	dto := UsuarioDTO{
		ID:     u.ID,
		Nome:   u.Nome,
		Email:  u.Email,
		Perfil: u.Perfil,
		Ativo:  u.Ativo,
	}

	return dto
}

// ToDetailDTO converte um modelo Usuario para UsuarioDetailDTO
func (u *Usuario) ToDetailDTO() UsuarioDetalheDTO {
	dto := UsuarioDetalheDTO{
		ID:                    u.ID,
		Nome:                  u.Nome,
		Email:                 u.Email,
		Perfil:                u.Perfil,
		Ativo:                 u.Ativo,
		ClinicaID:             u.ClinicaID,
		SupervisorID:          u.SupervisorID,
		DataInicioInatividade: u.DataInicioInatividade.Format("2006-01-02"),
		DataFimInatividade:    u.DataFimInatividade.Format("2006-01-02"),
		MotivoInatividade:     u.MotivoInatividade,
		CreatedAt:             u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:             u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return dto
}
