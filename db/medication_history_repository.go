package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
)

type MedicationHistoryRepository interface {
	CreateMedicationHistory(medicationHistory *models.MedicationHistory) (*models.MedicationHistory, error)
	UpdateMedicationHistory(hasMedicationBeenTaken bool, wasMedicationMissed string, medicationHistoryID uint, userID uint) error
	GetAllMedicationHistoryByUserID(userID uint) ([]models.MedicationHistory, error)
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

func (m *medicationHistoryRepo) UpdateMedicationHistory(hasMedicationBeenTaken bool, wasMedicationMissed string, medicationHistoryID uint, userID uint) error {
	err := m.DB.Model(&models.MedicationHistory{}).Select("has_medication_been_taken", "was_medication_missed").
		Where("user_id = ? AND id = ?", userID, medicationHistoryID).
		Updates(models.MedicationHistory{HasMedicationBeenTaken: hasMedicationBeenTaken, WasMedicationMissed: wasMedicationMissed}).Error
	if err != nil {
		return fmt.Errorf("could not update medication history: %v", err)
	}
	return nil
}

func (m *medicationHistoryRepo) GetAllMedicationHistoryByUserID(userID uint) ([]models.MedicationHistory, error) {
	var medicationHistories []models.MedicationHistory
	err := m.DB.Order("medication_time desc").Where("user_id = ?", userID).Find(&medicationHistories).Error
	if err != nil {
		return nil, fmt.Errorf("could not get medication history: %v", err)
	}
	return medicationHistories, nil
}
