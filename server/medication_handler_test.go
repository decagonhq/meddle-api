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

func TestCreateMedicationHandler(t *testing.T) {
	// generate a random user
	accToken, user := AuthorizeRoutes(t)

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

	// test cases
	testCases := []struct {
		name               string
		reqBody            interface{}
		medicationRequest  *models.MedicationRequest
		medicationResponse *models.MedicationResponse
		buildStubs         func(service *mocks.MockMedicationService, request *models.MedicationRequest, response *models.MedicationResponse)
		checkResponse      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success case",
			reqBody: gin.H{
				"name":                     "paracetamol",
				"dosage":                   2,
				"time_interval":            8,
				"medication_start_date":    "2012-04-23T18:25:43.511Z",
				"duration":                 7,
				"medication_prescribed_by": "Dr Tolu",
				"medication_stop_date":     "2012-04-23T18:25:43.511Z",
				"medication_start_time":    "2012-04-23T18:25:43.511Z",
				"purpose_of_medication":    "malaria treatment",
			},
			medicationRequest: &models.MedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2012-04-23T18:25:43.511Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     "2012-04-23T18:25:43.511Z",
				MedicationStartTime:    "2012-04-23T18:25:43.511Z",
				PurposeOfMedication:    "malaria treatment",
				UserID:                 user.ID,
			},
			medicationResponse: &models.MedicationResponse{
				ID:                     medication.ID,
				CreatedAt:              time.Unix(medication.CreatedAt, 0).String(),
				UpdatedAt:              time.Unix(medication.UpdatedAt, 0).String(),
				Name:                   medication.Name,
				Dosage:                 medication.Dosage,
				TimeInterval:           medication.TimeInterval,
				MedicationStartDate:    medication.MedicationStartDate.String(),
				Duration:               medication.Duration,
				MedicationPrescribedBy: medication.MedicationPrescribedBy,
				MedicationStopDate:     medication.MedicationStopDate.String(),
				MedicationStartTime:    medication.MedicationStartTime.String(),
				NextDosageTime:         medication.NextDosageTime.String(),
				PurposeOfMedication:    medication.PurposeOfMedication,
				UserID:                 user.ID,
			},
			buildStubs: func(service *mocks.MockMedicationService, request *models.MedicationRequest, response *models.MedicationResponse) {
				service.EXPECT().CreateMedication(request).Times(1).Return(response, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "bad request case",
			reqBody: gin.H{
				"name":                     "paracetamol",
				"dosage":                   2,
				"time_interval":            8,
				"medication_start_date":    "2012-04-23T18:25:43.511Z",
				"medication_prescribed_by": "Dr Tolu",
				"medication_stop_date":     "2012-04-23T18:25:43.511Z",
				"medication_start_time":    "2012-04-23T18:25:43.511Z",
				"purpose_of_medication":    "malaria treatment",
			},
			medicationRequest:  nil,
			medicationResponse: nil,
			buildStubs: func(service *mocks.MockMedicationService, request *models.MedicationRequest, response *models.MedicationResponse) {
				service.EXPECT().CreateMedication(request).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "bad request case due to date format",
			reqBody: gin.H{
				"name":                     "paracetamol",
				"dosage":                   2,
				"time_interval":            8,
				"medication_start_date":    "2012-04-23T18:25:43.511Z",
				"duration":                 7,
				"medication_prescribed_by": "Dr Tolu",
				"medication_stop_date":     "2012-04-23T18:25:43.511Z",
				"medication_start_time":    "2013-11-12",
				"purpose_of_medication":    "malaria treatment",
			},
			medicationRequest: &models.MedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2012-04-23T18:25:43.511Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     "2012-04-23T18:25:43.511Z",
				MedicationStartTime:    "2013-11-12",
				PurposeOfMedication:    "malaria treatment",
				UserID:                 user.ID,
			},
			medicationResponse: nil,
			buildStubs: func(service *mocks.MockMedicationService, request *models.MedicationRequest, response *models.MedicationResponse) {
				service.EXPECT().CreateMedication(request).Times(1).Return(nil, errors.ErrBadRequest)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "internal server error",
			reqBody: gin.H{
				"name":                     "paracetamol",
				"dosage":                   2,
				"time_interval":            8,
				"medication_start_date":    "2012-04-23T18:25:43.511Z",
				"duration":                 7,
				"medication_prescribed_by": "Dr Tolu",
				"medication_stop_date":     "2012-04-23T18:25:43.511Z",
				"medication_start_time":    "2012-04-23T18:25:43.511Z",
				"purpose_of_medication":    "malaria treatment",
			},
			medicationRequest: &models.MedicationRequest{
				Name:                   "paracetamol",
				Dosage:                 2,
				TimeInterval:           8,
				MedicationStartDate:    "2012-04-23T18:25:43.511Z",
				Duration:               7,
				MedicationPrescribedBy: "Dr Tolu",
				MedicationStopDate:     "2012-04-23T18:25:43.511Z",
				MedicationStartTime:    "2012-04-23T18:25:43.511Z",
				PurposeOfMedication:    "malaria treatment",
				UserID:                 user.ID,
			},
			medicationResponse: nil,
			buildStubs: func(service *mocks.MockMedicationService, request *models.MedicationRequest, response *models.MedicationResponse) {
				service.EXPECT().CreateMedication(request).Times(1).Return(nil, errors.ErrInternalServerError)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMedicationService := mocks.NewMockMedicationService(ctrl)
	mockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	testServer.handler.MedicationService = mockMedicationService
	testServer.handler.AuthRepository = mockAuthRepository

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthRepository.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
			mockAuthRepository.EXPECT().TokenInBlacklist(accToken).Return(false)

			tc.buildStubs(mockMedicationService, tc.medicationRequest, tc.medicationResponse)

			jsonFile, err := json.Marshal(tc.reqBody)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/api/v1/user/medications", strings.NewReader(string(jsonFile)))
			require.NoError(t, err)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accToken))

			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetNextMedicationHandler(t *testing.T) {
	// generate a random user
	accToken, user := AuthorizeRoutes(t)

	startDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
	stopDate, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")
	startTime, _ := time.Parse(time.RFC3339, "2013-10-21T13:28:06.419Z")

	medication := models.Medication{
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
	// test cases
	testCases := []struct {
		name               string
		medicationResponse []models.MedicationResponse
		buildStubs         func(service *mocks.MockMedicationService, userID uint, response []models.MedicationResponse)
		checkResponse      func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success case",
			medicationResponse: []models.MedicationResponse{
				{
					ID:                     medication.ID,
					CreatedAt:              time.Unix(medication.CreatedAt, 0).String(),
					UpdatedAt:              time.Unix(medication.UpdatedAt, 0).String(),
					Name:                   medication.Name,
					Dosage:                 medication.Dosage,
					TimeInterval:           medication.TimeInterval,
					MedicationStartDate:    medication.MedicationStartDate.String(),
					Duration:               medication.Duration,
					MedicationPrescribedBy: medication.MedicationPrescribedBy,
					MedicationStopDate:     medication.MedicationStopDate.String(),
					MedicationStartTime:    medication.MedicationStartTime.String(),
					NextDosageTime:         medication.NextDosageTime.String(),
					PurposeOfMedication:    medication.PurposeOfMedication,
					UserID:                 user.ID,
				},
			},
			buildStubs: func(service *mocks.MockMedicationService, request uint, response []models.MedicationResponse) {
				service.EXPECT().GetNextMedications(request).Times(1).Return(response, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:               "internal server error",
			medicationResponse: nil,
			buildStubs: func(service *mocks.MockMedicationService, request uint, response []models.MedicationResponse) {
				service.EXPECT().GetNextMedications(request).Times(1).Return(nil, errors.ErrInternalServerError)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMedicationService := mocks.NewMockMedicationService(ctrl)
	mockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	testServer.handler.MedicationService = mockMedicationService
	testServer.handler.AuthRepository = mockAuthRepository

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthRepository.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
			mockAuthRepository.EXPECT().TokenInBlacklist(accToken).Return(false)

			tc.buildStubs(mockMedicationService, user.ID, tc.medicationResponse)

			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/api/v1/user/medications/next", nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accToken))

			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
