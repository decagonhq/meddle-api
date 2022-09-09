package services

import (
	"net/http"
	"testing"

	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var mockRepository *mocks.MockAuthRepository
var testAuthService AuthService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockRepository = mocks.NewMockAuthRepository(ctrl)
	mailService := mocks.NewMockMailer(ctrl)
	pushNotification := mocks.NewMockPushNotifier(ctrl)
	testAuthService = NewAuthService(mockRepository, testConfig, mailService, pushNotification)

	mockMedicationRepository = mocks.NewMockMedicationRepository(ctrl)
	testMedicationService = NewMedicationService(mockMedicationRepository, mockMedicationHistoryRepository, testConfig)
	return func() {
		testAuthService = nil
		testMedicationService = nil
		defer ctrl.Finish()
	}
}

func Test_AuthLoginService(t *testing.T) {
	// arrange

	user := models.User{
		Model: models.Model{
			ID:        1,
			CreatedAt: 0,
			UpdatedAt: 0,
			DeletedAt: 0,
		},
		Name:           "name",
		PhoneNumber:    "1234567890",
		Email:          "email@gmail.com",
		Password:       "password",
		HashedPassword: "",
		IsEmailActive:  true,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)
	user.HashedPassword = string(hashedPassword)

	inactiveUser := user
	inactiveUser.IsEmailActive = false

	testCases := []struct {
		name          string
		input         models.LoginRequest
		dbOutput      *models.User
		dbError       error
		loginResponse *models.LoginResponse
		loginError    *errors.Error
	}{
		{
			name: "login successful case",
			input: models.LoginRequest{
				Email:    user.Email,
				Password: user.Password,
			},
			dbOutput: &user,
			dbError:  nil,
			loginResponse: &models.LoginResponse{
				UserResponse: models.UserResponse{
					ID:          user.ID,
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
				},
			},
			loginError: nil,
		},
		{
			name: "invalid email case",
			input: models.LoginRequest{
				Email:    "",
				Password: "password",
			},
			dbOutput:      nil,
			dbError:       gorm.ErrRecordNotFound,
			loginResponse: nil,
			loginError:    errors.New("invalid email", http.StatusUnprocessableEntity),
		},
		{
			name: "invalid password case",
			input: models.LoginRequest{
				Email:    user.Email,
				Password: "wrongpassword",
			},
			dbOutput:      &user,
			dbError:       nil,
			loginResponse: nil,
			loginError:    errors.ErrInvalidPassword,
		},
		{
			name: "inactive user",
			input: models.LoginRequest{
				Email:    inactiveUser.Email,
				Password: "password",
			},
			dbOutput:      &inactiveUser,
			dbError:       nil,
			loginResponse: nil,
			loginError:    errors.New("email not verified", http.StatusUnauthorized),
		},
		{
			name: "internal server error case",
			input: models.LoginRequest{
				Email:    user.Email,
				Password: user.Password,
			},
			dbOutput:      nil,
			dbError:       gorm.ErrInvalidDB,
			loginResponse: nil,
			loginError:    errors.ErrInternalServerError,
		},
	}
	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockRepository.EXPECT().FindUserByEmail(tc.input.Email).Times(1).Return(tc.dbOutput, tc.dbError)

			loginResponse, err := testAuthService.LoginUser(&tc.input)
			if tc.name != "login successful case" {
				require.Equal(t, tc.loginResponse, loginResponse)
				require.Equal(t, tc.loginError, err)
			}
		})
	}
}

func Test_DeleteUserByEmail(t *testing.T) {
	// arrange
	testCases := []struct {
		name                  string
		email                 string
		dbErrorOutput         error
		deleteUserErrorOutput *errors.Error
		buildStubs            func(repository *mocks.MockAuthRepository, email string, dbError error)
	}{
		{
			name:                  "delete user successful case",
			email:                 "sample@email.com",
			dbErrorOutput:         nil,
			deleteUserErrorOutput: nil,
			buildStubs: func(repository *mocks.MockAuthRepository, email string, dbError error) {
				repository.EXPECT().DeleteUserByEmail(email).Times(1).Return(dbError)
			},
		},
		{
			name:                  "delete user successful case",
			email:                 "sample@email.com",
			dbErrorOutput:         gorm.ErrInvalidDB,
			deleteUserErrorOutput: errors.ErrInternalServerError,
			buildStubs: func(repository *mocks.MockAuthRepository, email string, dbError error) {
				repository.EXPECT().DeleteUserByEmail(email).Times(1).Return(dbError)
			},
		},
	}

	teardown := setup(t)
	defer teardown()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockRepository, tc.email, tc.dbErrorOutput)
			err := testAuthService.DeleteUserByEmail(tc.email)

			require.Equal(t, tc.deleteUserErrorOutput, err)
		})
	}
}
