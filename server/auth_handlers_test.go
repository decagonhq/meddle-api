package server

import (
	"encoding/json"
	"fmt"
	"github.com/decagonhq/meddle-api/dto"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/services"
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

func Test_loginHandler(t *testing.T) {

	// generate a random user
	user, password := randomUser(t)
	// generate a token from a secret
	testSecret := "testSecret"
	token, err := services.GenerateToken(user.Email, &testSecret)
	require.NoError(t, err)
	//declare test cases
	testCases := []struct {
		name             string
		testRequest      dto.LoginRequest
		expectedResponse dto.LoginResponse
		buildStubs       func(service *mocks.MockAuthService, inputRequest dto.LoginRequest, secret string, response interface{})
	}{
		{
			name: "success",
			testRequest: dto.LoginRequest{
				Email:    user.Email,
				Password: password,
			},
			expectedResponse: dto.LoginResponse{
				UserResponse: dto.UserResponse{
					ID:          "",
					Name:        user.Name,
					PhoneNumber: user.PhoneNumber,
					Email:       user.Email,
				},
				AccessToken: *token,
			},
			buildStubs: func(service *mocks.MockAuthService, inputRequest dto.LoginRequest, secret string, response interface{}) {
				service.EXPECT().LoginUser(&inputRequest, secret).Times(1).Return(response, nil)
			},
		},
		//{
		//	name: "bad request",
		//	testRequest: dto.LoginRequest{
		//		Email:    "example@gmail.com",
		//		Password: "",
		//	},
		//	expectedResponse: dto.LoginResponse{
		//		UserResponse: dto.UserResponse{},
		//		AccessToken:  "",
		//	},
		//	buildStubs: func(service *mocks.MockAuthService, inputRequest dto.LoginRequest, secret string) {
		//		service.EXPECT().LoginUser(&inputRequest, secret).Times(0).Return(gomock.Any(), nil)
		//	},
		//},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := mocks.NewMockAuthService(ctrl)
			tc.buildStubs(mockService, tc.testRequest, testSecret, tc.expectedResponse)
			server := &Server{
				AuthService: mockService,
			}

			router := server.setupRouter()
			jsonFile, err := json.Marshal(tc.testRequest)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(string(jsonFile)))
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
			require.Contains(t, tc.expectedResponse.AccessToken, w.Body)
		})
	}
}

func randomUser(t *testing.T) (user models.User, password string) {
	password = RandomString(6)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user = models.User{
		Name:           RandomOwnerName(),
		HashedPassword: string(hashedPassword),
		PhoneNumber:    RandomOwnerName(),
		Email:          RandomEmail(),
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
