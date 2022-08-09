package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository

type MedicationRepository interface {
	CreateMedication(medication *models.Medication) (*models.Medication, error)
	GetMedicationDetail(id uint, userId uint) (*models.Medication, error)
	GetAllMedications(userID uint) ([]models.Medication, error)
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
		return nil, fmt.Errorf("could not create medication: %v", err)
	}
	return medication, nil
}

func (m *medicationRepo) GetMedicationDetail(id uint, userId uint) (*models.Medication, error) {
	var medication models.Medication
	err := m.DB.Where("id = ? AND user_id = ?", id, userId).First(&medication).Error
	if err != nil {
		return nil, fmt.Errorf("could not get medication: %v", err)
	}
	return &medication, nil
}

func (m *medicationRepo) GetAllMedications(userID uint) ([]models.Medication, error) {
	var medications []models.Medication
	err := m.DB.Where("user_id = ?", userID).Find(&medications).Error
	if err != nil {
		return nil, fmt.Errorf("could not get medications: %v", err)
	}
	return medications, nil
}

