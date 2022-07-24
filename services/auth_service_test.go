package services

import (
	"github.com/decagonhq/meddle-api/dto"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

var mockRepository *mocks.MockAuthRepository
var testLoginService AuthService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockRepository = mocks.NewMockAuthRepository(ctrl)
	testLoginService = NewAuthService(mockRepository)
	return func() {
		testLoginService = nil
		defer ctrl.Finish()
	}
}

func Test_AuthLoginService(t *testing.T) {
	// arrange

	user := &models.User{
		Model: models.Model{
			ID:        "id",
			CreatedAt: 0,
			UpdatedAt: 0,
			DeletedAt: 0,
		},
		Name:           "name",
		PhoneNumber:    "1234567890",
		Email:          "email@gmail.com",
		Password:       "password",
		HashedPassword: "",
		IsAgree:        true,
	}
	secret := "testJWTsecret"
	token, err := GenerateToken(user.Email, secret)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)
	user.HashedPassword = string(hashedPassword)
	testCases := []struct {
		name     string
		input    dto.LoginRequest
		secret   string
		dbOutput *models.User
		dbError  error
		output1  *dto.LoginResponse
		output2  *errors.Error
	}{
		{
			name: "login successful case",
			input: dto.LoginRequest{
				Email:    user.Email,
				Password: user.Password,
			},
			secret:   secret,
			dbOutput: user,
			dbError:  nil,
			output1: &dto.LoginResponse{
				UserResponse: dto.UserResponse{
					ID:          user.ID,
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
				},
				AccessToken: token,
			},
			output2: nil,
		},
		{
			name: "not found",
			input: dto.LoginRequest{
				Email:    "",
				Password: "password",
			},
			secret:   secret,
			dbOutput: nil,
			dbError:  gorm.ErrRecordNotFound,
			output1:  nil,
			output2:  errors.ErrNotFound,
		},
		{
			name: "invalid password",
			input: dto.LoginRequest{
				Email:    user.Email,
				Password: "wrongpassword",
			},
			secret:   secret,
			dbOutput: user,
			dbError:  nil,
			output1:  nil,
			output2:  errors.ErrInvalidPassword,
		},
		{
			name: "internal server error case",
			input: dto.LoginRequest{
				Email:    user.Email,
				Password: user.Password,
			},
			secret:   "",
			dbOutput: user,
			dbError:  nil,
			output1:  nil,
			output2:  errors.ErrInternalServerError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teardown := setup(t)
			defer teardown()
			mockRepository.EXPECT().FindUserByEmail(tc.input.Email).Times(1).Return(tc.dbOutput, tc.dbError)

			loginResponse, err := testLoginService.LoginUser(&tc.input, tc.secret)
			assert.Equal(t, tc.output1, loginResponse)
			assert.Equal(t, tc.output2, err)

			if tc.name == "login successful case" {
				assert.Equal(t, tc.output1.AccessToken, loginResponse.AccessToken)
			}
		})
	}

}
