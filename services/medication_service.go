package services

import (
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"log"
	"net/http"
	"time"
)

//go:generate mockgen -destination=../mocks/medication_mock.go -package=mocks github.com/decagonhq/meddle-api/services MedicationService

type MedicationService interface {
	CreateMedication(request *models.MedicationRequest) (*models.MedicationResponse, *errors.Error)
	GetNextMedications(userID uint) ([]models.MedicationResponse, *errors.Error)
}

// medicationService struct
type medicationService struct {
	Config         *config.Config
	medicationRepo db.MedicationRepository
}

// NewMedicationService instantiate an authService
func NewMedicationService(medicationRepo db.MedicationRepository, conf *config.Config) MedicationService {
	return &medicationService{
		Config:         conf,
		medicationRepo: medicationRepo,
	}
}

func (m *medicationService) CreateMedication(request *models.MedicationRequest) (*models.MedicationResponse, *errors.Error) {
	startDate, err := time.Parse(time.RFC3339, request.MedicationStartDate)
	if err != nil {
		return nil, errors.New("wrong date format", http.StatusBadRequest)
	}
	stopDate, err := time.Parse(time.RFC3339, request.MedicationStopDate)
	if err != nil {
		return nil, errors.New("wrong date format", http.StatusBadRequest)
	}
	startTime, err := time.Parse(time.RFC3339, request.MedicationStartTime)
	if err != nil {
		return nil, errors.New("wrong time format", http.StatusBadRequest)
	}

	medication := request.ReqToMedicationModel()
	medication.CreatedAt = time.Now().Unix()
	medication.UpdatedAt = time.Now().Unix()
	medication.MedicationStartDate = startDate
	medication.MedicationStopDate = stopDate
	medication.MedicationStartTime = startTime
	nextTime := medication.MedicationStartTime.Add(time.Hour * time.Duration(medication.TimeInterval))

	if nextTime.Day()-medication.MedicationStartTime.Day() <= 0 {
		medication.NextDosageTime = time.Date(nextTime.Year(), nextTime.Month(), nextTime.Day(), nextTime.Hour(), 0, 0, 0, time.UTC)
	} else {
		d := medication.MedicationStartTime
		medication.NextDosageTime = time.Date(d.Year(), d.Month(), d.Day()+1, 9, 0, 0, 0, time.UTC)
	}

	response, err := m.medicationRepo.CreateMedication(medication)
	if err != nil {
		log.Println(err)
		return nil, errors.ErrInternalServerError
	}
	return response.MedicationToResponse(), nil
}

func (m *medicationService) GetNextMedications(userID uint) ([]models.MedicationResponse, *errors.Error) {
	var nextMedicationResponses []models.MedicationResponse

	medications, err := m.medicationRepo.GetNextMedications(userID)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	for _, medication := range medications {
		nextMedicationResponses = append(nextMedicationResponses, *medication.MedicationToResponse())
	}
	return nextMedicationResponses, nil
}
