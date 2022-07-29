package services

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

var mockMedicationRepository *mocks.MockMedicationRepository
var testMedicationService MedicationService

func Test_CreateMedicationService(t *testing.T) {
	// arrange
	startDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
	stopDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
	startTime, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")

	medication := &models.Medication{
		Model: models.Model{
			ID:        0,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			DeletedAt: 0,
		},
		Name:                   "paracetamol",
		Dosage:                 2,
		TimeInterval:           8,
		MedicationStartDate:    startDate,
		Duration:               7,
		MedicationPrescribedBy: "Dr Tolu",
		MedicationStopDate:     stopDate,
		MedicationStartTime:    startTime,
		NextDosageTime:         startTime.Add(time.Hour * time.Duration(8)),
		PurposeOfMedication:    "malaria treatment",
	}
	testCases := []struct {
		name              string
		input             models.MedicationRequest
		dbInput           *models.Medication
		dbOutput          *models.Medication
		dbError           error
		createMedResponse *models.MedicationResponse
		createMedError    *errors.Error
		buildStubs        func(service *mocks.MockMedicationRepository, dbInput *models.Medication, dbOutput *models.Medication, dbError error)
	}{
		{
			name: "medication created successfully case",
			input: models.MedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2013-10-21T13:28:06.419Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     "2013-10-21T13:28:06.419Z",
				MedicationStartTime:    "2013-10-21T13:28:06.419Z",
				PurposeOfMedication:    "malaria treatment",
			},
			dbInput: &models.Medication{
				Name:                   medication.Name,
				Dosage:                 medication.Dosage,
				TimeInterval:           medication.TimeInterval,
				MedicationStartDate:    medication.MedicationStartDate,
				Duration:               medication.Duration,
				MedicationPrescribedBy: medication.MedicationPrescribedBy,
				MedicationStopDate:     medication.MedicationStopDate,
				MedicationStartTime:    medication.MedicationStartTime,
				PurposeOfMedication:    medication.PurposeOfMedication,
				IsMedicationDone:       medication.IsMedicationDone,
			},
			dbOutput: medication,
			dbError:  nil,
			createMedResponse: &models.MedicationResponse{
				ID:                     medication.ID,
				CreatedAt:              time.Unix(medication.CreatedAt, 0).String(),
				UpdatedAt:              time.Unix(medication.UpdatedAt, 0).String(),
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    medication.MedicationStartDate.String(),
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     medication.MedicationStopDate.String(),
				MedicationStartTime:    medication.MedicationStartTime.String(),
				NextDosageTime:         medication.NextDosageTime.String(),
				PurposeOfMedication:    "malaria treatment",
			},
			createMedError: nil,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, dbOutput *models.Medication, dbError error) {
				repository.EXPECT().CreateMedication(dbInput).Times(1).Return(dbOutput, dbError)
			},
		},
		{
			name: "bad request",
			input: models.MedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2013-10-21T13:28:06.419Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     "2013-10-21T13:28:06.419Z",
				MedicationStartTime:    "2013-11-12",
				PurposeOfMedication:    "malaria treatment",
			},
			dbInput:           nil,
			dbOutput:          nil,
			dbError:           nil,
			createMedResponse: nil,
			createMedError:    errors.New("wrong time format", http.StatusBadRequest),
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, dbOutput *models.Medication, dbError error) {
				repository.EXPECT().CreateMedication(dbInput).Times(0).Return(dbOutput, dbError)
			},
		},
		{
			name: "error creating medication due server error",
			input: models.MedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2013-10-21T13:28:06.419Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     "2013-10-21T13:28:06.419Z",
				MedicationStartTime:    "2013-10-21T13:28:06.419Z",
				PurposeOfMedication:    "malaria treatment",
			},
			dbInput: &models.Medication{
				Name:                   medication.Name,
				Dosage:                 medication.Dosage,
				TimeInterval:           medication.TimeInterval,
				MedicationStartDate:    medication.MedicationStartDate,
				Duration:               medication.Duration,
				MedicationPrescribedBy: medication.MedicationPrescribedBy,
				MedicationStopDate:     medication.MedicationStopDate,
				MedicationStartTime:    medication.MedicationStartTime,
				PurposeOfMedication:    medication.PurposeOfMedication,
				IsMedicationDone:       medication.IsMedicationDone,
			},
			dbOutput:          nil,
			dbError:           gorm.ErrInvalidDB,
			createMedResponse: nil,
			createMedError:    errors.ErrInternalServerError,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, dbOutput *models.Medication, dbError error) {
				repository.EXPECT().CreateMedication(dbInput).Times(1).Return(dbOutput, dbError)
			},
		},
	}
	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockMedicationRepository, medication, tc.dbOutput, tc.dbError)
			medicationResponse, err := testMedicationService.CreateMedication(&tc.input)

			require.Equal(t, tc.createMedResponse, medicationResponse)
			require.Equal(t, tc.createMedError, err)

		})
	}
}

func Test_GetAllMedicationsService(t *testing.T) {
	// arrange
	startDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
	stopDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
	startTime, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")

	medication := &models.Medication{
		Model: models.Model{
			ID:        1,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			DeletedAt: 0,
		},
		Name:                   "paracetamol",
		Dosage:                 2,
		TimeInterval:           8,
		MedicationStartDate:    startDate,
		Duration:               7,
		MedicationPrescribedBy: "Dr Tolu",
		MedicationStopDate:     stopDate,
		MedicationStartTime:    startTime,
		NextDosageTime:         startTime.Add(time.Hour * time.Duration(8)),
		PurposeOfMedication:    "malaria treatment",
	}
	testCases := []struct {
		name              string
		dbInput           uint
		dbOutput          []models.Medication
		dbError           error
		getAllMedResponse []models.MedicationResponse
		getAllMedError    *errors.Error
		buildStubs        func(repository *mocks.MockMedicationRepository, dbInput uint, dbOutput []models.Medication, dbError error)
	}{
		{
			name:    "getting all medications successful case",
			dbInput: 1,
			dbOutput: []models.Medication{
				{
					Model: models.Model{
						ID:        1,
						CreatedAt: time.Now().Unix(),
						UpdatedAt: time.Now().Unix(),
						DeletedAt: 0,
					},
					Name:                   "paracetamol",
					Dosage:                 2,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               7,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     stopDate,
					MedicationStartTime:    startTime,
					NextDosageTime:         startTime.Add(time.Hour * time.Duration(8)),
					PurposeOfMedication:    "malaria treatment",
					UserID:                 1,
				},
				{
					Model: models.Model{
						ID:        2,
						CreatedAt: time.Now().Unix(),
						UpdatedAt: time.Now().Unix(),
						DeletedAt: 0,
					},
					Name:                   "flagyl",
					Dosage:                 1,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               2,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     stopDate,
					MedicationStartTime:    startTime,
					NextDosageTime:         startTime.Add(time.Hour * time.Duration(8)),
					PurposeOfMedication:    "stomach pain",
					UserID:                 1,
				},
			},
			dbError: nil,
			getAllMedResponse: []models.MedicationResponse{
				{
					ID:                     medication.ID,
					CreatedAt:              time.Unix(medication.CreatedAt, 0).String(),
					UpdatedAt:              time.Unix(medication.UpdatedAt, 0).String(),
					Name:                   "paracetamol",
					Dosage:                 2,
					TimeInterval:           8,
					MedicationStartDate:    medication.MedicationStartDate.String(),
					Duration:               7,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     medication.MedicationStopDate.String(),
					MedicationStartTime:    medication.MedicationStartTime.String(),
					NextDosageTime:         medication.NextDosageTime.String(),
					PurposeOfMedication:    "malaria treatment",
					UserID:                 1,
				},
				{
					ID:                     medication.ID + 1,
					CreatedAt:              time.Unix(medication.CreatedAt, 0).String(),
					UpdatedAt:              time.Unix(medication.UpdatedAt, 0).String(),
					Name:                   "flagyl",
					Dosage:                 1,
					TimeInterval:           8,
					MedicationStartDate:    medication.MedicationStartDate.String(),
					Duration:               2,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     medication.MedicationStopDate.String(),
					MedicationStartTime:    medication.MedicationStartTime.String(),
					NextDosageTime:         medication.NextDosageTime.String(),
					PurposeOfMedication:    "stomach pain",
					UserID:                 1,
				},
			},
			getAllMedError: nil,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput uint, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetAllMedications(dbInput).Times(1).Return(dbOutput, dbError)
			},
		},
		{
			name:              "error creating medication due server error",
			dbInput:           1,
			dbOutput:          nil,
			dbError:           gorm.ErrInvalidDB,
			getAllMedResponse: nil,
			getAllMedError:    errors.ErrInternalServerError,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput uint, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetAllMedications(dbInput).Times(1).Return(dbOutput, dbError)
			},
		},
	}
	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockMedicationRepository, tc.dbInput, tc.dbOutput, tc.dbError)
			medicationResponse, err := testMedicationService.GetAllMedications(1)

			require.Equal(t, tc.getAllMedResponse, medicationResponse)
			require.Equal(t, tc.getAllMedError, err)

		})
	}

}
