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
	GetNextMedication(userID uint) (*models.Medication, error)
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

func (m *medicationRepo) GetNextMedication(userID uint) (*models.Medication, error) {
	var medication models.Medication
	err := m.DB.Where("user_id = ?", userID).Order("next_dosage_time ASC").First(&medication).Error
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("could not get next medication: %v", err)
	}
	return &medication, nil
}
