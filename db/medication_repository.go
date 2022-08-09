package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
	"time"
)

//go:generate mockgen -destination=../mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository

type MedicationRepository interface {
	CreateMedication(medication *models.Medication) (*models.Medication, error)
	GetNextMedications(userID uint) ([]models.Medication, error)
	UpdateMedicationDone(medication *models.Medication) error
	GetAllNextMedicationsToUpdate() ([]models.Medication, error)
	GetAllMedications(userID uint) ([]models.Medication, error)
	UpdateNextMedicationTime(medication *models.Medication, nextDosageTime time.Time) error
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

func (m *medicationRepo) GetNextMedications(userID uint) ([]models.Medication, error) {
	var medications []models.Medication
	err := m.DB.Where("user_id = ? AND next_dosage_time > ?", userID, time.Now().UTC()).Order("next_dosage_time ASC").Find(&medications).Error
	if err != nil {
		return nil, fmt.Errorf("could not get next medication: %v", err)
	}
	return medications, nil
}

func (m *medicationRepo) GetAllNextMedicationsToUpdate() ([]models.Medication, error) {
	var medications []models.Medication

	err := m.DB.Where("date_trunc('hour', next_dosage_time) = date_trunc('hour', now())").Where("is_medication_done = false").Find(&medications).Error
	if err != nil {
		return nil, fmt.Errorf("could not get next medication: %v", err)
	}
	return medications, nil
}

func (m *medicationRepo) UpdateMedicationDone(medication *models.Medication) error {
	err := m.DB.Model(&medication).Where("user_id = ?", medication.UserID).Update("is_medication_done", true).Error
	if err != nil {
		return fmt.Errorf("could not update medication: %v", err)
	}
	return nil
}

func (m *medicationRepo) UpdateNextMedicationTime(medication *models.Medication, nextDosageTime time.Time) error {
	err := m.DB.Model(&medication).Where("user_id = ?", medication.UserID).Update("next_dosage_time", nextDosageTime).Error
	if err != nil {
		return fmt.Errorf("could not update medication next time: %v", err)
	}
	return nil
}

func (m *medicationRepo) GetAllMedications(userID uint) ([]models.Medication, error) {
	var medications []models.Medication
	err := m.DB.Where("user_id = ?", userID).Find(&medications).Error
	if err != nil {
		return nil, fmt.Errorf("could not get medications: %v", err)
	}
	return medications, nil
}
