package services

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"time"
)

var testMedicationHistoryService MedicationHistoryService

func Test_UpdateMedicationHistoryService(t *testing.T) {
	// arrange

	testCases := []struct {
		name                   string
		reqInput               bool
		dbInput                string
		medicationID           uint
		userID                 uint
		dbError                error
		updateMedResponseError *errors.Error
		buildStubs             func(repository *mocks.MockMedicationHistoryRepository, dbInput string, medicationID uint, userID uint, dbError error)
	}{
		{
			name:                   "medication updates successfully case",
			reqInput:               true,
			dbInput:                "YES",
			dbError:                nil,
			updateMedResponseError: nil,
			buildStubs: func(repository *mocks.MockMedicationHistoryRepository, dbInput string, medicationID uint, userID uint, dbError error) {
				repository.EXPECT().UpdateMedicationHistory(true, dbInput, medicationID, userID).Times(1).Return(dbError)
			},
		},
		{
			name:                   "error updating medication due server error",
			reqInput:               true,
			dbInput:                "YES",
			dbError:                gorm.ErrInvalidDB,
			updateMedResponseError: errors.ErrInternalServerError,
			buildStubs: func(repository *mocks.MockMedicationHistoryRepository, dbInput string, medicationID uint, userID uint, dbError error) {
				repository.EXPECT().UpdateMedicationHistory(true, dbInput, medicationID, userID).Times(1).Return(dbError)
			},
		},
	}
	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockMedicationHistoryRepository, tc.dbInput, tc.medicationID, tc.userID, tc.dbError)
			err := testMedicationHistoryService.UpdateMedicationHistory(tc.reqInput, tc.medicationID, tc.userID)

			require.Equal(t, tc.updateMedResponseError, err)
		})
	}
}

func Test_GetAllMedicationHistoryByUserService(t *testing.T) {

	// arrange
	medicationTime, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")

	medicationHistory := models.MedicationHistory{
		MedicationName:         "paracetamol",
		MedicationID:           1,
		MedicationTime:         medicationTime,
		MedicationDosage:       2,
		UserID:                 1,
		HasMedicationBeenTaken: false,
		WasMedicationMissed:    "",
	}
	testCases := []struct {
		name              string
		dbInput           uint
		dbOutput          []models.MedicationHistory
		dbError           error
		getAllMedResponse []models.MedicationHistoryResponse
		getAllMedError    *errors.Error
		buildStubs        func(repository *mocks.MockMedicationHistoryRepository, dbInput uint, dbOutput []models.MedicationHistory, dbError error)
	}{
		{
			name:    "getting all medications successful case",
			dbInput: 1,
			dbOutput: []models.MedicationHistory{
				{
					MedicationID:           medicationHistory.MedicationID,
					MedicationName:         medicationHistory.MedicationName,
					MedicationDosage:       medicationHistory.MedicationDosage,
					MedicationTime:         medicationHistory.MedicationTime,
					HasMedicationBeenTaken: false,
					WasMedicationMissed:    "",
					UserID:                 1,
				},
				{
					Model: models.Model{
						ID: medicationHistory.ID + 1,
					},
					MedicationID:           medicationHistory.MedicationID,
					MedicationName:         medicationHistory.MedicationName,
					MedicationDosage:       medicationHistory.MedicationDosage,
					MedicationTime:         medicationHistory.MedicationTime,
					HasMedicationBeenTaken: false,
					WasMedicationMissed:    "",
					UserID:                 1,
				},
			},
			dbError: nil,
			getAllMedResponse: []models.MedicationHistoryResponse{
				{
					ID:                     medicationHistory.ID,
					CreatedAt:              time.Unix(medicationHistory.CreatedAt, 0).String(),
					UpdatedAt:              time.Unix(medicationHistory.UpdatedAt, 0).String(),
					MedicationID:           medicationHistory.MedicationID,
					MedicationName:         medicationHistory.MedicationName,
					MedicationDosage:       medicationHistory.MedicationDosage,
					MedicationTime:         medicationHistory.MedicationTime.UTC().String(),
					HasMedicationBeenTaken: false,
					UserID:                 1,
				},
				{
					ID:                     medicationHistory.ID + 1,
					CreatedAt:              time.Unix(medicationHistory.CreatedAt, 0).String(),
					UpdatedAt:              time.Unix(medicationHistory.UpdatedAt, 0).String(),
					MedicationID:           medicationHistory.MedicationID,
					MedicationName:         medicationHistory.MedicationName,
					MedicationDosage:       medicationHistory.MedicationDosage,
					MedicationTime:         medicationHistory.MedicationTime.UTC().String(),
					HasMedicationBeenTaken: false,
					UserID:                 1,
				},
			},
			getAllMedError: nil,
			buildStubs: func(repository *mocks.MockMedicationHistoryRepository, dbInput uint, dbOutput []models.MedicationHistory, dbError error) {
				repository.EXPECT().GetAllMedicationHistoryByUserID(dbInput).Times(1).Return(dbOutput, dbError)
			},
		},
		{
			name:              "error creating medication due server error",
			dbInput:           1,
			dbOutput:          nil,
			dbError:           gorm.ErrInvalidDB,
			getAllMedResponse: nil,
			getAllMedError:    errors.ErrInternalServerError,
			buildStubs: func(repository *mocks.MockMedicationHistoryRepository, dbInput uint, dbOutput []models.MedicationHistory, dbError error) {
				repository.EXPECT().GetAllMedicationHistoryByUserID(dbInput).Times(1).Return(dbOutput, dbError)
			},
		},
	}
	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockMedicationHistoryRepository, tc.dbInput, tc.dbOutput, tc.dbError)
			medicationResponse, err := testMedicationHistoryService.GetAllMedicationHistoryByUser(1)

			require.Equal(t, tc.getAllMedResponse, medicationResponse)
			require.Equal(t, tc.getAllMedError, err)
		})
	}

}
