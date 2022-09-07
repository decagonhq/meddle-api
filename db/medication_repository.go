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
	GetMedicationDetail(id uint, userId uint) (*models.Medication, error)
	GetAllMedications(userID uint) ([]models.Medication, error)
	UpdateNextMedicationTime(medication *models.Medication, nextDosageTime time.Time) error
	UpdateMedication(medication *models.Medication, medicationID uint, userID uint) error
	FindMedication(medicationName, by, purpose string, duration int, dosage int,) (*[]models.Medication, error)
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

func (m *medicationRepo) UpdateMedication(medication *models.Medication, medicationID uint, userID uint) error {
	err := m.DB.Model(&models.Medication{}).
		Where("user_id = ? AND id = ?", userID, medicationID).
		Updates(medication).Error
	if err != nil {
		return fmt.Errorf("could not update medication: %v", err)
	}
	return nil
}

func (m *medicationRepo) FindMedication(medicationName, by, purpose string,  duration int, dosage int) (*[]models.Medication, error) {
	var medications *[]models.Medication
	 err := m.DB.Where("name = ?", medicationName).Or("dosage = ?", dosage).Or("duration = ?",duration).Or("medication_prescribed_by = ?",by).Or("purpose_of_medication = ?",purpose).Find(&medications).Error
	 if err != nil{
		 return nil, err
	 }
	 return medications, nil
}

