package server

import (
	"encoding/json"
	"fmt"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_UpdateMedicationHistoryHandler(t *testing.T) {

	// generate a random user
	accToken, user := AuthorizeTestUser(t)

	// test cases
	testCases := []struct {
		name                string
		reqBody             interface{}
		routeParam          string
		medicationHistoryID uint
		reqBodyValue        bool
		errorResponse       *errors.Error
		buildStubs          func(service *mocks.MockMedicationHistoryService, reqBodyValue bool, medicationHistoryID uint, userID uint, errorResponse *errors.Error)
		checkResponse       func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success case",
			reqBody: gin.H{
				"has_medication_been_taken": true,
			},
			reqBodyValue:        true,
			medicationHistoryID: 1,
			routeParam:          "1",
			buildStubs: func(service *mocks.MockMedicationHistoryService, reqBodyValue bool, medicationID uint, userID uint, errorResponse *errors.Error) {
				service.EXPECT().UpdateMedicationHistory(reqBodyValue, medicationID, userID).Times(1).Return(errorResponse)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "internal server error",
			reqBody: gin.H{
				"has_medication_been_taken": true,
			},
			reqBodyValue:        true,
			medicationHistoryID: 1,
			routeParam:          "1",
			errorResponse:       errors.ErrInternalServerError,
			buildStubs: func(service *mocks.MockMedicationHistoryService, reqBodyValue bool, medicationID uint, userID uint, errorResponse *errors.Error) {
				service.EXPECT().UpdateMedicationHistory(reqBodyValue, medicationID, userID).Times(1).Return(errorResponse)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:       "bad request from route param",
			routeParam: "a",
			buildStubs: func(service *mocks.MockMedicationHistoryService, reqBodyValue bool, medicationID uint, userID uint, errorResponse *errors.Error) {
				service.EXPECT().UpdateMedicationHistory(reqBodyValue, medicationID, userID).Times(0).Return(errorResponse)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMedicationHistoryService := mocks.NewMockMedicationHistoryService(ctrl)
	mockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	testServer.handler.MedicationHistoryService = mockMedicationHistoryService
	testServer.handler.AuthRepository = mockAuthRepository

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthRepository.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
			mockAuthRepository.EXPECT().TokenInBlacklist(accToken).Return(false)

			tc.buildStubs(mockMedicationHistoryService, tc.reqBodyValue, tc.medicationHistoryID, user.ID, tc.errorResponse)

			jsonFile, err := json.Marshal(tc.reqBody)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/user/medication-history/%v", tc.routeParam), strings.NewReader(string(jsonFile)))
			require.NoError(t, err)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accToken))

			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAllMedicationHistoryByUserHandler(t *testing.T) {

	// generate a random user
	accToken, user := AuthorizeTestUser(t)
	medicationTime, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")

	medicationHistory := models.MedicationHistory{
		MedicationName:         "paracetamol",
		MedicationID:           1,
		MedicationTime:         medicationTime,
		MedicationDosage:       2,
		UserID:                 user.ID,
		HasMedicationBeenTaken: false,
		WasMedicationMissed:    "",
	}

	// test cases
	testCases := []struct {
		name               string
		medicationResponse []models.MedicationHistoryResponse
		buildStubs         func(service *mocks.MockMedicationHistoryService, userID uint, response []models.MedicationHistoryResponse)
		checkCodeResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success case",
			medicationResponse: []models.MedicationHistoryResponse{
				{
					ID:                     medicationHistory.ID,
					CreatedAt:              time.Unix(medicationHistory.CreatedAt, 0).String(),
					UpdatedAt:              time.Unix(medicationHistory.UpdatedAt, 0).String(),
					MedicationID:           medicationHistory.MedicationID,
					MedicationName:         medicationHistory.MedicationName,
					MedicationDosage:       medicationHistory.MedicationDosage,
					MedicationTime:         medicationHistory.MedicationTime.String(),
					UserID:                 user.ID,
					HasMedicationBeenTaken: false,
					WasMedicationMissed:    "",
				},
				{
					ID:                     medicationHistory.ID + 1,
					CreatedAt:              time.Unix(medicationHistory.CreatedAt, 0).String(),
					UpdatedAt:              time.Unix(medicationHistory.UpdatedAt, 0).String(),
					MedicationID:           medicationHistory.MedicationID,
					MedicationName:         medicationHistory.MedicationName,
					MedicationDosage:       medicationHistory.MedicationDosage,
					MedicationTime:         medicationHistory.MedicationTime.String(),
					UserID:                 user.ID,
					HasMedicationBeenTaken: true,
					WasMedicationMissed:    "YES",
				},
			},
			buildStubs: func(service *mocks.MockMedicationHistoryService, request uint, response []models.MedicationHistoryResponse) {
				service.EXPECT().GetAllMedicationHistoryByUser(request).Times(1).Return(response, nil)
			},
			checkCodeResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:               "internal server error",
			medicationResponse: nil,
			buildStubs: func(service *mocks.MockMedicationHistoryService, request uint, response []models.MedicationHistoryResponse) {
				service.EXPECT().GetAllMedicationHistoryByUser(request).Times(1).Return(nil, errors.ErrInternalServerError)
			},
			checkCodeResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMedicationHistoryService := mocks.NewMockMedicationHistoryService(ctrl)
	mockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	testServer.handler.MedicationHistoryService = mockMedicationHistoryService
	testServer.handler.AuthRepository = mockAuthRepository

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthRepository.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
			mockAuthRepository.EXPECT().TokenInBlacklist(accToken).Return(false)

			tc.buildStubs(mockMedicationHistoryService, user.ID, tc.medicationResponse)

			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/user/medication-history", nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accToken))

			testServer.router.ServeHTTP(recorder, req)
			tc.checkCodeResponse(t, recorder)
		})
	}
}
