package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
)

type MedicationHistoryRepository interface {
	CreateMedicationHistory(medicationHistory *models.MedicationHistory) (*models.MedicationHistory, error)
	UpdateMedicationHistory(medication *models.MedicationHistory, medicationHistoryID uint, userID uint) error
	GetAllMedicationHistory(userID uint) ([]models.MedicationHistory, error)
}

type medicationHistoryRepo struct {
	DB *gorm.DB
}

func NewMedicationHistoryRepo(db *GormDB) MedicationHistoryRepository {
	return &medicationHistoryRepo{db.DB}
}

func (m *medicationHistoryRepo) CreateMedicationHistory(medicationHistory *models.MedicationHistory) (*models.MedicationHistory, error) {
	err := m.DB.Create(medicationHistory).Error
	if err != nil {
		return nil, fmt.Errorf("could not create medication: %v", err)
	}
	return medicationHistory, nil
}
func (m *medicationHistoryRepo) UpdateMedicationHistory(medication *models.MedicationHistory, medicationHistoryID uint, userID uint) error {
	return nil
}
func (m *medicationHistoryRepo) GetAllMedicationHistory(userID uint) ([]models.MedicationHistory, error) {
	return nil, nil
}
