package services

import (
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
var testLoginService AuthService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockRepository = mocks.NewMockAuthRepository(ctrl)
	testLoginService = NewAuthService(mockRepository, testConfig)

	mockMedicationRepository = mocks.NewMockMedicationRepository(ctrl)
	testMedicationService = NewMedicationService(mockMedicationRepository, testConfig)
	return func() {
		testLoginService = nil
		testMedicationService = nil
		defer ctrl.Finish()
	}
}

func Test_AuthLoginService(t *testing.T) {
	// arrange

	user := &models.User{
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
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)
	user.HashedPassword = string(hashedPassword)
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
			dbOutput: user,
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
			name: "not found",
			input: models.LoginRequest{
				Email:    "",
				Password: "password",
			},
			dbOutput:      nil,
			dbError:       gorm.ErrRecordNotFound,
			loginResponse: nil,
			loginError:    errors.ErrNotFound,
		},
		{
			name: "invalid password",
			input: models.LoginRequest{
				Email:    user.Email,
				Password: "wrongpassword",
			},
			dbOutput:      user,
			dbError:       nil,
			loginResponse: nil,
			loginError:    errors.ErrInvalidPassword,
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

			loginResponse, err := testLoginService.LoginUser(&tc.input)
			if tc.name != "login successful case" {
				require.Equal(t, tc.loginResponse, loginResponse)
				require.Equal(t, tc.loginError, err)
			} else {
				//require.NotZero(t, loginResponse.AccessToken)
				//require.Equal(t, tc.loginError, err)
			}

		})
	}
}
