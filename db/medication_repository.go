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
	err := m.DB.Where("user_id = ? AND next_dosage_time > ?", userID, time.Now().UTC()).Order("next_dosage_time ASC").Find(&medications).Error
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
		log.Println(fmt.Errorf("could not get next medications while running update next dosage cron job: %v", err))
	}
	for _, medication := range medications {
		timeSumation := medication.NextDosageTime.Add(time.Hour * time.Duration(medication.TimeInterval))
		setNextDosageTime := SetNextDosageTime(timeSumation, medication.NextDosageTime)

		if medication.IsMedicationDone == false && medication.NextDosageTime.Unix() < medication.MedicationStopDate.Unix() {
			m.DB.Model(&medication).Where("user_id = ?", medication.UserID).Update("next_dosage_time", setNextDosageTime)
		}
	}
}

func UpdateNextMedicationCronJob(repo MedicationRepository) {
	_, presentMinute, presentSecond := time.Now().UTC().Clock()
	waitTime := time.Duration(60-presentMinute)*time.Minute + time.Duration(60-presentSecond)*time.Second
	s := gocron.NewScheduler(time.UTC)
	time.Sleep(waitTime)
	s.Every(1).Hour().Do(func() {
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

func SetNextDosageTime(t1, t2 time.Time) time.Time {
	if t1.Day()-t2.Day() <= 0 {
		return time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), 0, 0, 0, time.UTC)
	}
	return time.Date(t2.Year(), t2.Month(), t2.Day()+1, 9, 0, 0, 0, time.UTC)
}
