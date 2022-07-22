package server

import (
	"encoding/json"
	"fmt"
	"github.com/decagonhq/meddle-api/dto"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_LoginHandler(t *testing.T) {
	testSecret := testServer.handler.Config.JWTSecret
	// generate a random user
	user, password := randomUser(t)

	// test cases
	testCases := []struct {
		name          string
		reqBody       interface{}
		loginRequest  *dto.LoginRequest
		loginResponse *dto.LoginResponse
		buildStubs    func(service *mocks.MockAuthService, request *dto.LoginRequest, secret string, response *dto.LoginResponse)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "success case",
			reqBody: gin.H{
				"email":    user.Email,
				"password": password,
			},
			loginRequest: &dto.LoginRequest{
				Email:    user.Email,
				Password: password,
			},
			loginResponse: &dto.LoginResponse{
				UserResponse: dto.UserResponse{
					ID:          user.ID,
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
				},
				AccessToken: "",
			},
			buildStubs: func(service *mocks.MockAuthService, request *dto.LoginRequest, secret string, response *dto.LoginResponse) {
				service.EXPECT().LoginUser(request, secret).Times(1).Return(response, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "invalid password case",
			reqBody: gin.H{
				"email":    user.Email,
				"password": "invalid password",
			},
			loginRequest: &dto.LoginRequest{
				Email:    user.Email,
				Password: "invalid password",
			},
			loginResponse: &dto.LoginResponse{
				UserResponse: dto.UserResponse{
					ID:          user.ID,
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
				},
				AccessToken: "",
			},
			buildStubs: func(service *mocks.MockAuthService, request *dto.LoginRequest, secret string, response *dto.LoginResponse) {
				service.EXPECT().LoginUser(request, secret).Times(1).Return(nil, errors.ErrUnauthorized)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "bad request case",
			reqBody: gin.H{
				"email":    user.Email,
				"password": "",
			},
			loginRequest:  nil,
			loginResponse: nil,
			buildStubs: func(service *mocks.MockAuthService, request *dto.LoginRequest, secret string, response *dto.LoginResponse) {
				service.EXPECT().LoginUser(request, secret).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "not found case",
			reqBody: gin.H{
				"email":    "user@email.com",
				"password": password,
			},
			loginRequest: &dto.LoginRequest{
				Email:    "user@email.com",
				Password: password,
			},
			loginResponse: &dto.LoginResponse{
				UserResponse: dto.UserResponse{
					ID:          user.ID,
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
				},
				AccessToken: "",
			},
			buildStubs: func(service *mocks.MockAuthService, request *dto.LoginRequest, secret string, response *dto.LoginResponse) {
				service.EXPECT().LoginUser(request, secret).Times(1).Return(nil, errors.ErrNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "internal server error case",
			reqBody: gin.H{
				"email":    user.Email,
				"password": password,
			},
			loginRequest: &dto.LoginRequest{
				Email:    user.Email,
				Password: password,
			},
			loginResponse: &dto.LoginResponse{
				UserResponse: dto.UserResponse{
					ID:          user.ID,
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
				},
				AccessToken: "",
			},
			buildStubs: func(service *mocks.MockAuthService, request *dto.LoginRequest, secret string, response *dto.LoginResponse) {
				service.EXPECT().LoginUser(request, secret).Times(1).Return(nil, errors.ErrInternalServerError)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockAuthService(ctrl)
			testServer.handler.AuthService = mockService
			tc.buildStubs(mockService, tc.loginRequest, testSecret, tc.loginResponse)

			jsonFile, err := json.Marshal(tc.reqBody)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(string(jsonFile)))
			require.NoError(t, err)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user models.User, password string) {
	password = RandomString(6)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user = models.User{
		Model: models.Model{
			ID:        RandomString(6),
			CreatedAt: RandomInt(1, 100),
			UpdatedAt: RandomInt(1, 100),
			DeletedAt: RandomInt(1, 100),
		},
		Name:           RandomOwnerName(),
		HashedPassword: string(hashedPassword),
		PhoneNumber:    RandomOwnerName(),
		Email:          RandomEmail(),
		IsAgree:        true,
	}
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const alphabet = "abcdefghhijklmnopqrstuvwxyz"

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString geerates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwnerName generates random account owner names for testing
func RandomOwnerName() string {
	return RandomString(6)
}

// RandomEmail generates a random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
