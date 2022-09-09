package services

import (
	"fmt"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"log"
	"net/http"
	"testing"
	"time"
)

var mockMedicationRepository *mocks.MockMedicationRepository
var mockMedicationHistoryRepository *mocks.MockMedicationHistoryRepository
var testMedicationService MedicationService

func Test_CreateMedicationService(t *testing.T) {
	// arrange
	startDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
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
		MedicationStopDate:     startTime.AddDate(0, 0, 7),
		MedicationStartTime:    startTime,
		NextDosageTime:         time.Date(startTime.Add(time.Hour*time.Duration(8)).Year(), startTime.Add(time.Hour*time.Duration(8)).Month(), startTime.Add(time.Hour*time.Duration(8)).Day(), startTime.Add(time.Hour*time.Duration(8)).Hour(), 0, 0, 0, time.UTC),
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
	startTime, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")

	medication := &models.Medication{
		Name:                   "paracetamol",
		Dosage:                 2,
		TimeInterval:           8,
		MedicationStartDate:    startDate,
		Duration:               7,
		MedicationPrescribedBy: "Dr Tolu",
		MedicationStopDate:     startTime.AddDate(0, 0, 7),
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
					Name:                   "paracetamol",
					Dosage:                 2,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               7,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     startTime.AddDate(0, 0, 7),
					MedicationStartTime:    startTime,
					NextDosageTime:         startTime.Add(time.Hour * time.Duration(8)),
					PurposeOfMedication:    "malaria treatment",
					UserID:                 1,
				},
				{
					Model: models.Model{
						ID: medication.ID + 1,
					},
					Name:                   "flagyl",
					Dosage:                 1,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               2,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     startTime.AddDate(0, 0, 2),
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
					MedicationStopDate:     medication.MedicationStartDate.AddDate(0, 0, 7).String(),
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
					MedicationStopDate:     medication.MedicationStartDate.AddDate(0, 0, 2).String(),
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

func Test_GetNextMedicationService(t *testing.T) {
	// arrange
	startDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
	startTime, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")

	medication := &models.Medication{
		Name:                   "paracetamol",
		Dosage:                 2,
		TimeInterval:           8,
		MedicationStartDate:    startDate,
		Duration:               7,
		MedicationPrescribedBy: "Dr Tolu",
		MedicationStopDate:     startTime.AddDate(0, 0, 7),
		MedicationStartTime:    startTime,
		NextDosageTime:         startTime.Add(time.Hour * time.Duration(8)),
		PurposeOfMedication:    "malaria treatment",
	}
	testCases := []struct {
		name               string
		dbInput            uint
		dbOutput           []models.Medication
		dbError            error
		getNextMedResponse []models.MedicationResponse
		getNextMedError    *errors.Error
		buildStubs         func(repository *mocks.MockMedicationRepository, dbInput uint, dbOutput []models.Medication, dbError error)
	}{
		{
			name:    "getting next medications successful case",
			dbInput: 1,
			dbOutput: []models.Medication{
				{
					Name:                   "paracetamol",
					Dosage:                 2,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               7,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     startTime.AddDate(0, 0, 7),
					MedicationStartTime:    startTime,
					NextDosageTime:         startTime.Add(time.Hour * time.Duration(8)),
					PurposeOfMedication:    "malaria treatment",
					UserID:                 1,
				},
			},
			dbError: nil,
			getNextMedResponse: []models.MedicationResponse{
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
			},
			getNextMedError: nil,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput uint, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetNextMedications(dbInput).Times(1).Return(dbOutput, dbError)
			},
		},
		{
			name:               "error getting next medications due server error",
			dbInput:            1,
			dbOutput:           nil,
			dbError:            gorm.ErrInvalidDB,
			getNextMedResponse: nil,
			getNextMedError:    errors.ErrInternalServerError,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput uint, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetNextMedications(dbInput).Times(1).Return(dbOutput, dbError)
			},
		},
	}
	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockMedicationRepository, tc.dbInput, tc.dbOutput, tc.dbError)
			medicationResponse, err := testMedicationService.GetNextMedications(1)

			require.Equal(t, tc.getNextMedResponse, medicationResponse)
			require.Equal(t, tc.getNextMedError, err)

		})
	}

}

func Test_CronUpdateMedicationForNextTime(t *testing.T) {
	startDate := time.Now().UTC()
	stopDate := startDate.AddDate(0, 0, 7)
	startTime := startDate

	testCases := []struct {
		name          string
		dbInput       *models.Medication
		dbOutput      []models.Medication
		dbError       error
		buildStubs    func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, timeInput time.Time, dbOutput []models.Medication, dbError error)
		checkResponse func(t *testing.T, cronJobError error)
	}{
		{
			name: "updating medication's next time successful case",
			dbInput: &models.Medication{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    startDate,
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     stopDate,
				MedicationStartTime:    startTime,
				NextDosageTime:         startTime,
				PurposeOfMedication:    "malaria treatment",
				UserID:                 1,
			},
			dbOutput: []models.Medication{
				{
					Name:                   "paracetamol",
					Dosage:                 2,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               7,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     stopDate,
					MedicationStartTime:    startTime,
					NextDosageTime:         startTime,
					PurposeOfMedication:    "malaria treatment",
					UserID:                 1,
				},
			},
			dbError: nil,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, timeInput time.Time, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetAllNextMedicationsToUpdate().Times(1).Return(dbOutput, dbError)
				repository.EXPECT().UpdateNextMedicationTime(dbInput, timeInput).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, cronJobError error) {
				require.Nil(t, cronJobError)
			},
		},
		{
			name: "updating medication's is done successful case",
			dbInput: &models.Medication{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    startDate,
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     stopDate.AddDate(0, 0, -7),
				MedicationStartTime:    startTime,
				NextDosageTime:         startTime,
				PurposeOfMedication:    "malaria treatment",
				UserID:                 1,
			},
			dbOutput: []models.Medication{
				{
					Name:                   "paracetamol",
					Dosage:                 2,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               7,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     stopDate.AddDate(0, 0, -7),
					MedicationStartTime:    startTime,
					NextDosageTime:         startTime,
					PurposeOfMedication:    "malaria treatment",
					UserID:                 1,
				},
			},
			dbError: nil,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, timeInput time.Time, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetAllNextMedicationsToUpdate().Times(1).Return(dbOutput, dbError)
				repository.EXPECT().UpdateMedicationDone(dbInput).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, cronJobError error) {
				log.Println(cronJobError)
				require.Nil(t, cronJobError)
			},
		},
		{
			name: "error getting next medications to update case",
			dbInput: &models.Medication{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    startDate,
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     stopDate,
				MedicationStartTime:    startTime,
				NextDosageTime:         startTime,
				PurposeOfMedication:    "malaria treatment",
				UserID:                 1,
			},
			dbOutput: nil,
			dbError:  fmt.Errorf("could not get next medication: %v", gorm.ErrInvalidDB),
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, timeInput time.Time, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetAllNextMedicationsToUpdate().Times(1).Return(dbOutput, dbError)
			},
			checkResponse: func(t *testing.T, cronJobError error) {
				log.Println(cronJobError)
				require.EqualError(t, cronJobError, fmt.Sprint("could not get next medications while running update next dosage cron job"))
			},
		},
		{
			name: "error updating medication's next time case",
			dbInput: &models.Medication{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    startDate,
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     stopDate,
				MedicationStartTime:    startTime,
				NextDosageTime:         startTime,
				PurposeOfMedication:    "malaria treatment",
				UserID:                 1,
			},
			dbOutput: []models.Medication{
				{
					Name:                   "paracetamol",
					Dosage:                 2,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               7,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     stopDate,
					MedicationStartTime:    startTime,
					NextDosageTime:         startTime,
					PurposeOfMedication:    "malaria treatment",
					UserID:                 1,
				},
			},
			dbError: nil,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, timeInput time.Time, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetAllNextMedicationsToUpdate().Times(1).Return(dbOutput, dbError)
				repository.EXPECT().UpdateNextMedicationTime(dbInput, timeInput).Times(1).Return(fmt.Errorf("could not update medication: %v", gorm.ErrInvalidDB))
			},
			checkResponse: func(t *testing.T, cronJobError error) {
				require.EqualError(t, cronJobError, fmt.Sprint("could not update next medication time while running update next dosage cron job"))
			},
		},
		{
			name: "error updating medication's is done fail case",
			dbInput: &models.Medication{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    startDate,
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     stopDate.AddDate(0, 0, -7),
				MedicationStartTime:    startTime,
				NextDosageTime:         startTime,
				PurposeOfMedication:    "malaria treatment",
				UserID:                 1,
			},
			dbOutput: []models.Medication{
				{
					Name:                   "paracetamol",
					Dosage:                 2,
					TimeInterval:           8,
					MedicationStartDate:    startDate,
					Duration:               7,
					MedicationPrescribedBy: "Dr Tolu",
					MedicationStopDate:     stopDate.AddDate(0, 0, -7),
					MedicationStartTime:    startTime,
					NextDosageTime:         startTime,
					PurposeOfMedication:    "malaria treatment",
					UserID:                 1,
				},
			},
			dbError: nil,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, timeInput time.Time, dbOutput []models.Medication, dbError error) {
				repository.EXPECT().GetAllNextMedicationsToUpdate().Times(1).Return(dbOutput, dbError)
				repository.EXPECT().UpdateMedicationDone(dbInput).Times(1).Return(fmt.Errorf("could not update medication: %v", gorm.ErrInvalidDB))
			},
			checkResponse: func(t *testing.T, cronJobError error) {
				require.EqualError(t, cronJobError, fmt.Sprint("could not update is medication done while running update next dosage cron job"))
			},
		},
	}
	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			timeSumation := tc.dbInput.NextDosageTime.Add(time.Hour * time.Duration(tc.dbInput.TimeInterval))
			nextDosageTime := GetNextDosageTime(timeSumation, tc.dbInput.NextDosageTime)
			tc.buildStubs(mockMedicationRepository, tc.dbInput, nextDosageTime, tc.dbOutput, tc.dbError)
			err := testMedicationService.CronUpdateMedicationForNextTime()

			tc.checkResponse(t, err)

		})
	}

}

func Test_UpdateMedicationService(t *testing.T) {
	// arrange
	startDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
	startTime, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")

	medication := &models.Medication{
		Name:                   "paracetamol",
		Dosage:                 2,
		TimeInterval:           8,
		MedicationStartDate:    startDate,
		Duration:               7,
		MedicationPrescribedBy: "Dr Tolu",
		MedicationStopDate:     startTime.AddDate(0, 0, 7),
		MedicationStartTime:    startTime,
		NextDosageTime:         time.Date(startTime.Add(time.Hour*time.Duration(8)).Year(), startTime.Add(time.Hour*time.Duration(8)).Month(), startTime.Add(time.Hour*time.Duration(8)).Day(), startTime.Add(time.Hour*time.Duration(8)).Hour(), 0, 0, 0, time.UTC),
		PurposeOfMedication:    "malaria treatment",
	}
	testCases := []struct {
		name                   string
		input                  models.UpdateMedicationRequest
		dbInput                *models.Medication
		medicationID           uint
		userID                 uint
		dbError                error
		updateMedResponseError *errors.Error
		buildStubs             func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, medicationID uint, userID uint, dbError error)
	}{
		{
			name: "medication updates successfully case",
			input: models.UpdateMedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2013-10-21T13:28:06.419Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
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
				MedicationStopDate:     medication.MedicationStartTime.AddDate(0, 0, 7),
				MedicationStartTime:    medication.MedicationStartTime,
				PurposeOfMedication:    medication.PurposeOfMedication,
				NextDosageTime:         medication.NextDosageTime,
			},
			dbError:                nil,
			updateMedResponseError: nil,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, medicationID uint, userID uint, dbError error) {
				repository.EXPECT().UpdateMedication(dbInput, medicationID, userID).Times(1).Return(dbError)
			},
		},
		{
			name: "bad request",
			input: models.UpdateMedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2013-10-21T13:28:06.419Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStartTime:    "2013-11-12",
				PurposeOfMedication:    "malaria treatment",
			},
			dbInput:                nil,
			dbError:                nil,
			updateMedResponseError: errors.New("wrong time format", http.StatusBadRequest),
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, medicationID uint, userID uint, dbError error) {
				repository.EXPECT().UpdateMedication(dbInput, medicationID, userID).Times(0).Return(dbError)
			},
		},
		{
			name: "error updating medication due server error",
			input: models.UpdateMedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2013-10-21T13:28:06.419Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
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
				MedicationStopDate:     medication.MedicationStartTime.AddDate(0, 0, 7),
				MedicationStartTime:    medication.MedicationStartTime,
				PurposeOfMedication:    medication.PurposeOfMedication,
				NextDosageTime:         medication.NextDosageTime,
			},
			dbError:                gorm.ErrInvalidDB,
			updateMedResponseError: errors.ErrInternalServerError,
			buildStubs: func(repository *mocks.MockMedicationRepository, dbInput *models.Medication, medicationID uint, userID uint, dbError error) {
				repository.EXPECT().UpdateMedication(dbInput, medicationID, userID).Times(1).Return(dbError)
			},
		},
	}
	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockMedicationRepository, tc.dbInput, tc.medicationID, tc.userID, tc.dbError)
			err := testMedicationService.UpdateMedication(&tc.input, tc.medicationID, tc.userID)

			require.Equal(t, tc.updateMedResponseError, err)
		})
	}
}
