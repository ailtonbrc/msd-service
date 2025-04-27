package service

import (
	"clinica_server/internal/auth"
	"clinica_server/internal/models"
	"clinica_server/internal/repository"
	"clinica_server/internal/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

// AuthService gerencia a autenticação de usuários
type AuthService struct {
    userRepo   repository.UsuarioRepository
    jwtService auth.JWTService
}

// NewAuthService cria um novo serviço de autenticação
func NewAuthService(userRepo repository.UsuarioRepository, jwtService auth.JWTService) *AuthService {
    return &AuthService{
        userRepo:   userRepo,
        jwtService: jwtService,
    }
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
    Usuario      models.UsuarioDTO `json:"usuario"`
    AccessToken  string            `json:"access_token"`
    RefreshToken string            `json:"refresh_token"`
    ExpiresIn    int               `json:"expires_in"`
}

// getPermissionsForRole retorna as permissões para um determinado perfil
func (s *AuthService) getPermissionsForRole(role string) []string {
    // Mapeamento de perfis para permissões
    rolePermissions := map[string][]string{
        "admin": {
            "pacientes:view", "pacientes:create", "pacientes:update", "pacientes:delete",
            "usuarios:view", "usuarios:create", "usuarios:update", "usuarios:delete",
            // Adicione outras permissões conforme necessário
        },
        "medico": {
            "pacientes:view", "pacientes:create", "pacientes:update",
            // Adicione outras permissões conforme necessário
        },
        "atendente": {
            "pacientes:view", "pacientes:create",
            // Adicione outras permissões conforme necessário
        },
        "usuario": {
            "pacientes:view",
            // Adicione outras permissões conforme necessário
        },
    }

    // Retorna as permissões para o perfil ou uma lista vazia se o perfil não existir
    if permissions, exists := rolePermissions[role]; exists {
        return permissions
    }
    return []string{}
}

// Login autentica um usuário e retorna tokens JWT
func (s *AuthService) Login(email, senha string) (*LoginResponse, error) {
    // Buscar usuário pelo email usando o repositório
    user, err := s.userRepo.BuscaPorEmail(email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("usuário não encontrado")
        }
        return nil, err
    }

    // Verificar se o usuário está ativo
    if !user.Ativo {
        return nil, errors.New("usuário inativo")
    }

    // Verifica se a senha informada corresponde ao hash armazenado
    if !utils.CheckPasswordHash(senha, user.Senha) {
        return nil, errors.New("senha incorreta")
    }

    // Definir permissões com base no perfil do usuário
    roles := []string{user.Perfil}
    permissions := s.getPermissionsForRole(user.Perfil)
    scopes := []string{}      // Preencher conforme necessário

    // Gerar tokens usando o JWTService
    accessToken, err := s.jwtService.GenerateToken(
        user.ID,
        user.Email,
        user.Email,
        roles,
        permissions,
        scopes,
        time.Hour, // Duração do token de acesso
    )
    if err != nil {
        return nil, err
    }

    refreshToken, err := s.jwtService.GenerateToken(
        user.ID,
        user.Email,
        user.Email,
        roles,
        permissions,
        scopes,
        time.Hour*24*7, // Duração do token de refresh (7 dias)
    )
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Usuario:      user.ToDTO(),
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    60, // 1 hora em minutos
    }, nil
}

// RefreshToken renova o token de acesso usando um token de refresh
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
    // Validar token de refresh
    claims, err := s.jwtService.ValidateToken(refreshToken)
    if err != nil {
        return nil, err
    }

    // Buscar usuário pelo email usando o repositório
    user, err := s.userRepo.BuscaPorEmail(claims.Email)
    if err != nil {
        return nil, err
    }

    // Verificar se o usuário está ativo
    if !user.Ativo {
        return nil, errors.New("usuário inativo")
    }

    // Atualizar permissões com base no perfil atual do usuário
    // Isso garante que as permissões estejam atualizadas mesmo se o perfil do usuário mudar
    roles := []string{user.Perfil}
    permissions := s.getPermissionsForRole(user.Perfil)

    // Gerar novo token usando o JWTService
    newAccessToken, err := s.jwtService.GenerateToken(
        user.ID,
        user.Email,
        user.Email,
        roles,
        permissions,
        claims.Scopes,
        time.Hour, // Duração do token de acesso
    )
    if err != nil {
        return nil, err
    }

    // Gerar novo token de refresh
    newRefreshToken, err := s.jwtService.GenerateToken(
        user.ID,
        user.Email,
        user.Email,
        roles,
        permissions,
        claims.Scopes,
        time.Hour*24*7, // Duração do token de refresh (7 dias)
    )
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Usuario:      user.ToDTO(),
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken,
        ExpiresIn:    60, // 1 hora em minutos
    }, nil
}

// GetUserByID busca um usuário pelo ID
func (s *AuthService) GetUserByID(userID uint) (*models.Usuario, error) {
    // Buscar usuário pelo ID usando o repositório
    user, err := s.userRepo.BuscaPorID(userID)
    if err != nil {
        return nil, err
    }
    return user, nil
}