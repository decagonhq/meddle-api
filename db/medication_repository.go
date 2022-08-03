package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"github.com/go-co-op/gocron"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"time"
)

//go:generate mockgen -destination=../mocks/medication_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db MedicationRepository

type MedicationRepository interface {
	CreateMedication(medication *models.Medication) (*models.Medication, error)
	GetNextMedication(userID uint) ([]models.Medication, error)
	UpdateNextMedicationTime()
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

func (m *medicationRepo) GetNextMedication(userID uint) ([]models.Medication, error) {
	var medication []models.Medication
	err := m.DB.Where("user_id = ? AND next_dosage_time > ?", userID, time.Now().UTC()).Order("next_dosage_time ASC").Limit(10).Find(&medication).Error
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("could not get next medication: %v", err)
	}
	return medication, nil
}

// remove minute and seconds use only hour for time
// 1 Fetch all medications if the next time for dosage => present time...
// 2 Update the next time for each medication to the sum of (next dosage time and the time interval) only if
// the sum is doesn't enter the next day, otherwise set the next dosage time to 8am/9am the next day
// 3 Check if the next dosage time is != dosage stop time
//fetch in batches...

func (m *medicationRepo) GetAllNextMedicationsToUpdate() ([]models.Medication, error) {
	var medications []models.Medication
	err := m.DB.Where("next_dosage_time = ?", time.Now().UTC()).Find(&medications).Error
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("could not get next medication: %v", err)
	}
	return medications, nil
}

func (m *medicationRepo) UpdateNextMedicationTime() {
	medications, err := m.GetAllNextMedicationsToUpdate()
	if err != nil {
		log.Println(err)
		return
	}
	for _, medication := range medications {
		timeSumation := medication.NextDosageTime.Add(time.Hour * time.Duration(medication.TimeInterval))
		diff := timeSumation.Day() - medication.NextDosageTime.Day()
		// check if the timeSumation has not entered the next day
		if diff == 0 {
			medication.NextDosageTime = timeSumation
			m.DB.Model(&medication).Where("user_id = ?", medication.UserID).Update("next_dosage_time", medication.NextDosageTime)
		} else {
			// next day at 8am or 9am
			d := medication.NextDosageTime
			medication.NextDosageTime = time.Date(d.Year(), d.Month(), d.Day()+1, 9, 0, 0, 0, time.UTC)
			m.DB.Model(&medication).Where("user_id = ?", medication.UserID).Update("next_dosage_time", medication.NextDosageTime)
		}
	}
}

func UpdateNextMedicationCronJob(repo MedicationRepository) {

	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Seconds().Do(func() {
		repo.UpdateNextMedicationTime()
	})
	s.StartBlocking()
}

func RunCronJob(repo MedicationRepository) {
	c := cron.New()
	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	c.AddFunc("@every 5m", func() { repo.UpdateNextMedicationTime() })
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
