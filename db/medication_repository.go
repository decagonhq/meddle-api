package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository

type MedicationRepository interface {
	CreateMedication(user *models.Medication) (*models.Medication, error)
}

type medicationRepo struct {
	DB *gorm.DB
}

func NewMedicationRepo(db *GormDB) MedicationRepository {
	return &medicationRepo{db.DB}
}

func (m *medicationRepo) CreateMedication(medication *models.Medication) (*models.Medication, error) {
	err := m.DB.Create(medication).Error
	if err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}
	return medication, nil
}

