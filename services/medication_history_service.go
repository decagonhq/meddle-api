package services

import (
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"log"
)

//go:generate mockgen -destination=../mocks/medication_history_mock.go -package=mocks github.com/decagonhq/meddle-api/services MedicationHistoryService

type MedicationHistoryService interface {
	UpdateMedicationHistory(hasMedicationBeenTaken bool, medicationHistoryID uint, userID uint) *errors.Error
	GetAllMedicationHistoryByUser(userID uint) ([]models.MedicationHistoryResponse, *errors.Error)
}

// medicationHistoryService struct
type medicationHistoryService struct {
	Config                *config.Config
	medicationHistoryRepo db.MedicationHistoryRepository
}

// NewMedicationHistoryService instantiate an authService
func NewMedicationHistoryService(medicationHistoryRepo db.MedicationHistoryRepository, conf *config.Config) MedicationHistoryService {
	return &medicationHistoryService{
		Config:                conf,
		medicationHistoryRepo: medicationHistoryRepo,
	}
}

func (m *medicationHistoryService) UpdateMedicationHistory(hasMedicationBeenTaken bool, medicationHistoryID uint, userID uint) *errors.Error {
	var wasMedicationMissed string
	if hasMedicationBeenTaken == true {
		wasMedicationMissed = "NO"
	} else {
		wasMedicationMissed = "YES"
	}
	err := m.medicationHistoryRepo.UpdateMedicationHistory(hasMedicationBeenTaken, wasMedicationMissed, medicationHistoryID, userID)
	if err != nil {
		log.Printf("error updating medication history: %v", err)
		return errors.ErrInternalServerError
	}
	return nil
}

func (m *medicationHistoryService) GetAllMedicationHistoryByUser(userID uint) ([]models.MedicationHistoryResponse, *errors.Error) {
	var medicationHistoryResponses []models.MedicationHistoryResponse

	medicationHistories, err := m.medicationHistoryRepo.GetAllMedicationHistoryByUserID(userID)
	if err != nil {
		log.Printf("error getting all medication history of user %v : %v", userID, err)
		return nil, errors.ErrInternalServerError
	}

	for _, medicationHistory := range medicationHistories {
		medicationHistoryResponses = append(medicationHistoryResponses, *medicationHistory.MedicationHistoryToResponse())
	}
	return medicationHistoryResponses, nil
}
