package service

import (
	"clinica_server/internal/models"
	"clinica_server/internal/repository"
	"clinica_server/internal/utils"
	"errors"
)

// UsuarioService gerencia operações relacionadas a usuários
type UsuarioService struct {
	userRepo repository.UsuarioRepository
}

// NewUsuarioService cria um novo serviço de usuários
func NewUsuarioService(repo repository.UsuarioRepository) *UsuarioService {
	return &UsuarioService{
		userRepo: repo,
	}
}

// GetUsuarios retorna uma lista paginada de usuários
func (s *UsuarioService) BuscaUsuarios(pagination *utils.Pagination) (*models.UsuarioListDTO, error) {
	usuarios, err := s.userRepo.BuscaTodos(pagination)
	if err != nil {
		return nil, err
	}

	// Converter para DTOs
	usuarioDTOs := make([]models.UsuarioDTO, 0, len(usuarios))
	for _, usuario := range usuarios {
		usuarioDTOs = append(usuarioDTOs, usuario.ToDTO())
	}

	return &models.UsuarioListDTO{
		Usuarios:   usuarioDTOs,
		Pagination: models.ToPaginationDTO(pagination),
	}, nil
}

// GetUsuarioByID busca um usuário pelo ID
func (s *UsuarioService) BuscaUsuarioPorID(id uint) (*models.UsuarioDetalheDTO, error) {
	usuario, err := s.userRepo.BuscaPorID(id)
	if err != nil {
		return nil, err
	}
	if usuario == nil {
		return nil, utils.ErrNotFound
	}

	// Converter para DTO
	usuarioDetailDTO := usuario.ToDetailDTO()
	return &usuarioDetailDTO, nil
}

// CreateUsuario cria um novo usuário
func (s *UsuarioService) CreateUsuario(req models.CreateUserRequest) (*models.UsuarioDTO, error) {

	// Validação de Dados para Criar o Usuário

	// Verificar se o email já existe
	user, err := s.userRepo.BuscaPorEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("email já está em uso")
	}

	// Cria o Hash da senha
	passwordHash, err := utils.HashPassword(req.Senha)
	if err != nil {
		return nil, err
	}

	// Criar usuário
	usuario := models.Usuario{
		Nome:   req.Nome,
		Email:  req.Email,
		Senha:  passwordHash,
		Perfil: req.Perfil,
		Ativo:  true, // Por padrão, usuários são criados ativos
	}

	if err := s.userRepo.Create(&usuario); err != nil {
		return nil, err
	}

	// Buscar usuário completo com relacionamentos
	completeUsuario, err := s.userRepo.BuscaPorID(usuario.ID)
	if err != nil {
		return nil, err
	}

	// Converter para DTO
	usuarioDTO := completeUsuario.ToDTO()
	return &usuarioDTO, nil
}

// UpdateUsuario atualiza um usuário existente
func (s *UsuarioService) UpdateUsuario(id uint, req models.UpdateUserRequest) (*models.UsuarioDTO, error) {
	// Validação de Dados para Criar o Usuário

	// Verifica se o usuário existe
	usuario, err := s.userRepo.BuscaPorID(id)
	if err != nil {
		return nil, err
	}
	if usuario == nil {
		return nil, utils.ErrNotFound
	}

	// Verificar se o email já está em uso por outro usuário (se fornecido)
	if req.Email != nil && *req.Email != usuario.Email {
		exists, err := s.userRepo.VerificaOutroUsuarioPorEmail(*req.Email, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email já está em uso")
		}
	}

	// Atualizar campos
	if *req.Nome != "" {
		usuario.Nome = *req.Nome
	}
	if req.Email != nil {
		usuario.Email = *req.Email
	}
	if req.Perfil != nil {
		usuario.Perfil = *req.Perfil
	}
	if req.Ativo != nil {
		usuario.Ativo = *req.Ativo
	}
	if req.ClinicaID != nil {
		usuario.ClinicaID = *req.ClinicaID
	}
	if req.SupervisorID != nil {
		usuario.SupervisorID = *req.SupervisorID
	}
	if req.DataInicioInatividade != nil {
		usuario.DataInicioInatividade = *req.DataInicioInatividade
	}
	if req.DataFimInatividade != nil {
		usuario.DataFimInatividade = *req.DataFimInatividade
	}
	if req.MotivoInatividade != nil {
		usuario.MotivoInatividade = *req.MotivoInatividade
	}

	// Salvar alterações
	if err := s.userRepo.Update(usuario); err != nil {
		return nil, err
	}

	// Buscar usuário completo com relacionamentos
	completeUsuario, err := s.userRepo.BuscaPorID(usuario.ID)
	if err != nil {
		return nil, err
	}

	// Converter para DTO
	usuarioDTO := completeUsuario.ToDTO()
	return &usuarioDTO, nil
}

// ChangePassword altera a senha de um usuário
func (s *UsuarioService) TrocaSenha(id uint, senhaAtual, novaSenha string, isAdmin bool) error {
	// Validar dados
	req := models.ChangePasswordRequest{
		CurrentPassword: senhaAtual,
		NewPassword:     novaSenha,
	}

	// Verificar se o usuário existe
	user, err := s.userRepo.BuscaPorID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("usuário não encontrado")
	}

	// Verificar nova senha
	if len(req.NewPassword) < 6 {
		return errors.New("nova senha deve ter pelo menos 6 caracteres")
	}

	// Buscar usuário
	usuario, err := s.userRepo.BuscaPorID(id)
	if err != nil {
		return err
	}
	if usuario == nil {
		return utils.ErrNotFound
	}

	// TODO Precisa ver como q vai ser a questão de Perfil
	//Se não for admin, verificar a senha atual
	if !isAdmin && !utils.CheckPasswordHash(senhaAtual, usuario.Senha) {
		return utils.ErrInvalidCredentials
	}

	// Hash da nova senha
	passwordHash, err := utils.HashPassword(novaSenha)
	if err != nil {
		return err
	}

	// Atualizar senha
	usuario.Senha = passwordHash
	return s.userRepo.Update(usuario)
}

// DeleteUsuario exclui um usuário (soft delete)
func (s *UsuarioService) DeleteUsuario(id uint) error {
	// Verificar se o usuário existe
	usuario, err := s.userRepo.BuscaPorID(id)
	if err != nil {
		return err
	}
	if usuario == nil {
		return utils.ErrNotFound
	}

	// Excluir usuário
	return s.userRepo.Delete(id)
}
