package server

import (
	"bytes"
	"encoding/json"
	"github.com/decagonhq/meddle-api/services/jwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetPassword(t *testing.T) {
	user := models.User{}
	email := "toluwasethomas1@gmail.com"
	token, _ := jwt.GenerateToken(email, testServer.handler.Config.JWTSecret)
	newReq := &models.ResetPassword{
		Password:        "12345678",
		ConfirmPassword: "12345678",
	}
	user.HashedPassword, _ = services.GenerateHashPassword(newReq.Password)

	cases := []struct {
		Name            string
		Request         *models.ResetPassword
		ExpectedCode    int
		ExpectedMessage string
		ExpectedError   string
		mockDB          func(ctrl *mocks.MockAuthRepository)
		checkResponse   func(recorder *httptest.ResponseRecorder)
	}{
		{
			Name:            "Test Reset Password",
			Request:         newReq,
			ExpectedCode:    http.StatusCreated,
			ExpectedMessage: "Reset successful, Login with your new password to continue",
			ExpectedError:   "",
			mockDB: func(ctrl *mocks.MockAuthRepository) {
				ctrl.EXPECT().IsTokenInBlacklist(token).Return(nil).AnyTimes()
				ctrl.EXPECT().UpdatePassword(gomock.Any(), email).Return(nil).AnyTimes()
				ctrl.EXPECT().AddToBlackList(gomock.Any()).Return(nil).AnyTimes()
			},
		},
		{
			Name:            "Test Supply with password not equal",
			Request:         &models.ResetPassword{Password: newReq.Password, ConfirmPassword: "hhhhhhhh"},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "",
			ExpectedError:   "password does not match",
			mockDB:          func(ctrl *mocks.MockAuthRepository) {},
		},
		{
			Name:            "Test Supply with short password",
			Request:         &models.ResetPassword{Password: "abcd", ConfirmPassword: "abcd"},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "",
			ExpectedError:   "wrong password length",
			mockDB:          func(ctrl *mocks.MockAuthRepository) {},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
	mail := mocks.NewMockMailer(ctrl)
	pushNotifier := mocks.NewMockPushNotifier(ctrl)
	authService := services.NewAuthService(mockAuthRepo, testServer.handler.Config, mail, pushNotifier)
	testServer.handler.AuthService = authService
	testServer.handler.AuthRepository = mockAuthRepo

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			c.mockDB(mockAuthRepo)
			data, err := json.Marshal(c.Request)
			require.NoError(t, err)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/password/reset/"+token, bytes.NewReader(data))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			assert.Equal(t, recorder.Code, c.ExpectedCode)
		})
	}
}
