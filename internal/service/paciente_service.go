package service

import (
	"clinica_server/internal/models"
	"errors"

	"gorm.io/gorm"
)

type PacienteService struct {
	db *gorm.DB
}

func NewPacienteService(db *gorm.DB) *PacienteService {
	return &PacienteService{db: db}
}

func (s *PacienteService) GetPacientes() ([]models.Paciente, error) {
	var pacientes []models.Paciente
	err := s.db.Find(&pacientes).Error
	return pacientes, err
}

func (s *PacienteService) GetPacienteByID(id uint) (*models.Paciente, error) {
	var paciente models.Paciente
	if err := s.db.First(&paciente, id).Error; err != nil {
		return nil, err
	}
	return &paciente, nil
}

func (s *PacienteService) CreatePaciente(req models.CreatePacienteRequest) (*models.Paciente, error) {
	paciente := models.Paciente{
		Nome:           req.Nome,
		DataNascimento: req.DataNascimento,
		Genero:         req.Genero,
		Diagnostico:    req.Diagnostico,
		Telefone:       req.Telefone,
		Email:          req.Email,
		Endereco:       req.Endereco,
	}
	err := s.db.Create(&paciente).Error
	return &paciente, err
}

func (s *PacienteService) UpdatePaciente(id uint, req models.UpdatePacienteRequest) (*models.Paciente, error) {
	paciente, err := s.GetPacienteByID(id)
	if err != nil {
		return nil, err
	}
	s.db.Model(paciente).Updates(req)
	return paciente, nil
}

func (s *PacienteService) DeletePaciente(id uint) error {
	result := s.db.Delete(&models.Paciente{}, id)
	if result.RowsAffected == 0 {
		return errors.New("paciente não encontrado")
	}
	return result.Error
}