package service

import (
	"clinica_server/config"
	"clinica_server/internal/models"
	"clinica_server/internal/utils"
	"errors"

	"gorm.io/gorm"
)

// AuthService gerencia a autenticação de usuários
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService cria um novo serviço de autenticação
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
	Usuario      models.UsuarioDTO `json:"usuario"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int               `json:"expires_in"`
}

// Login autentica um usuário e retorna tokens JWT
func (s *AuthService) Login(email, senha string) (*LoginResponse, error) {
	var user models.Usuario

	// Buscar usuário pelo username
	result := s.db.Where("LOWER(email) = LOWER(?)", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, result.Error
	}

	// Verificar se o usuário está ativo
	if !user.Ativo {
		return nil, errors.New("usuário inativo")
	}

	// senhash, err := utils.HashPassword(password)
	// fmt.Printf("Senha hash: %s\n", senhash)

	// // Verificar senha
	// fmt.Printf("Senha: %s, SenhaBD: %s", password, user.PasswordHash)

	// Verifica se a senha informada corresponde ao hash armazenado
	if !utils.CheckPasswordHash(senha, user.Senha) {
		return nil, errors.New("senha incorreta")
	}

	// Gerar tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Perfil, s.cfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Usuario:      user.ToDTO(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.cfg.JWT.AccessTokenExp.Minutes()),
	}, nil
}

// RefreshToken renova o token de acesso usando um token de refresh
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// Validar token de refresh
	claims, err := utils.ValidateToken(refreshToken, s.cfg)
	if err != nil {
		return nil, err
	}

	// Buscar usuário
	var user models.Usuario
	result := s.db.Where("LOWER(email) = LOWER(?)", claims.Subject).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// Verificar se o usuário está ativo
	if !user.Ativo {
		return nil, errors.New("usuário inativo")
	}

	// Gerar novo token de acesso
	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Perfil, s.cfg)
	if err != nil {
		return nil, err
	}

	// Gerar novo token de refresh
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Usuario:      user.ToDTO(),
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int(s.cfg.JWT.AccessTokenExp.Minutes()),
	}, nil
}

// GetUserByID busca um usuário pelo ID
func (s *AuthService) GetUserByID(userID uint) (*models.Usuario, error) {
	var user models.Usuario
	result := s.db.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
