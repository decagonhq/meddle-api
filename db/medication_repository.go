package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
	"log"
	"time"
)

//go:generate mockgen -destination=../mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository

type MedicationRepository interface {
	CreateMedication(medication *models.Medication) (*models.Medication, error)
	GetNextMedications(userID uint) ([]models.Medication, error)
	UpdateNextMedicationTime()
	GetAllNextMedicationsToUpdate() ([]models.Medication, error)
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

func (m *medicationRepo) GetNextMedications(userID uint) ([]models.Medication, error) {
	var medications []models.Medication
	err := m.DB.Where("user_id = ? AND next_dosage_time > ?", userID, time.Now().UTC()).Order("next_dosage_time ASC").Limit(10).Find(&medications).Error
	if err != nil {
		return nil, fmt.Errorf("could not get next medication: %v", err)
	}
	return medications, nil
}

func (m *medicationRepo) GetAllNextMedicationsToUpdate() ([]models.Medication, error) {
	var medications []models.Medication

	err := m.DB.Where("(SELECT date_trunc('minute', next_dosage_time)) = (SELECT date_trunc('minute', now()))").Find(&medications).Error
	if err != nil {
		return nil, fmt.Errorf("could not get next medication: %v", err)
	}
	return medications, nil
}

func (m *medicationRepo) UpdateNextMedicationTime() {
	medications, err := m.GetAllNextMedicationsToUpdate()
	if err != nil {
		log.Println(err)
	}
	for _, medication := range medications {
		timeSumation := medication.NextDosageTime.Add(time.Hour * time.Duration(medication.TimeInterval))
		diff := timeSumation.Day() - medication.NextDosageTime.Day()

		if medication.NextDosageTime != medication.MedicationStopDate && medication.IsMedicationDone == false && medication.NextDosageTime.Unix() < medication.MedicationStopDate.Unix() {
			if diff == 0 {
				m.DB.Model(&medication).Where("user_id = ?", medication.UserID).Update("next_dosage_time", timeSumation)
			} else {
				medication.NextDosageTime = SetNextDosageTime(medication.NextDosageTime)
				m.DB.Model(&medication).Where("user_id = ?", medication.UserID).Update("next_dosage_time", medication.NextDosageTime)
			}
		}
	}
}

func UpdateNextMedicationCronJob(repo MedicationRepository) {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(func() {
		repo.UpdateNextMedicationTime()
	})
	s.StartBlocking()
}

func (m *medicationRepo) GetAllMedications(userID uint) ([]models.Medication, error) {
	var medications []models.Medication
	err := m.DB.Where("user_id = ?", userID).Find(&medications).Error
	if err != nil {
		return nil, fmt.Errorf("could not get medications: %v", err)
	}
	return medications, nil
}

func SetNextDosageTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+1, 9, 0, 0, 0, time.UTC)
}
