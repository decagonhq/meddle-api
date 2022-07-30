package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

//go:generate mockgen -destination=../mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository

type MedicationRepository interface {
	CreateMedication(medication *models.Medication) (*models.Medication, error)
	GetMedicationById(userId uint, medId string) (*models.Medication, *errors.Error)
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

func (m *medicationRepo) GetMedicationById(userId uint, medId string) (*models.Medication, *errors.Error) {
	var medication models.Medication
	err := m.DB.Where("id = ? AND user_id= ?", medId, userId).First(&medication).Error
	if err != nil {
		return nil, errors.New("medication not found", http.StatusNotFound)
	}
	return &medication, nil
}