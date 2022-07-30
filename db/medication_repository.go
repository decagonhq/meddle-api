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
	GetMedicationById(userId string) (*models.User, error)
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

func (m *medicationRepo) GetMedicationById(userId string) (*models.User, error) {
	var user models.User
	err := m.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("could not get medication: %v", err)
	}
	return &user, nil
}