package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
	"log"
)

//go:generate mockgen -destination=../mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository

type MedicationRepository interface {
	CreateMedication(medication *models.Medication) (*models.Medication, error)
	GetMedicationById(userId uint) (*models.Medication, error)
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
		log.Println(err)
		return nil, fmt.Errorf("could not create medication: %v", err)
	}
	return medication, nil
}

func (m *medicationRepo) GetMedicationById(userId uint) (*models.Medication, error) {
	var medication models.Medication
	err := m.DB.Where("id = ?", userId).First(&medication).Error
	if err != nil {
		return nil, fmt.Errorf("could not get medication: %v", err)
	}
	return &medication, nil
}