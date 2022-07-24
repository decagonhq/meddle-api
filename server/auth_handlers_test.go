package server

import (
	"bytes"
	"encoding/json"
	"github.com/decagonhq/meddle-api/db/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignup(t *testing.T) {
	newReq := &models.User{
		Name:        "Tolu",
		PhoneNumber: "+2348163608141",
		Email:       "toluwase@gmail.com",
		Password:    "12345678",
	}
	newResp := &models.User{
		Name:        newReq.Name,
		PhoneNumber: newReq.PhoneNumber,
		Email:       newReq.Email,
	}
	cases := []struct {
		Name            string
		Request         *models.User
		ExpectedCode    int
		ExpectedMessage string
		ExpectedError   string
		buildStubs      func(ctrl *mocks.MockAuthService)
		checkResponse   func(recorder *httptest.ResponseRecorder)
	}{
		{
			Name:            "Test Signup with correct details",
			Request:         newReq,
			ExpectedCode:    http.StatusCreated,
			ExpectedMessage: "user created successfully",
			ExpectedError:   "",
		},
		{
			Name:            "Test Signup with no email",
			Request:         &models.User{Name: "Tolu", PhoneNumber: "08141636082"},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "",
			ExpectedError:   "Email is invalid: toluwase.tt.com",
		},
		{
			Name:            "Test Signup with invalid fields",
			Request:         &models.User{Name: "Tolu", PhoneNumber: "08141", Email: "tolut.a"},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "",
			ExpectedError:   "Email is invalid: toluwase.tt.com",
		},
		{
			Name:            "Test Signup with duplicate email address",
			Request:         &models.User{Name: "Tolu", PhoneNumber: "08141", Email: "toluwase@gmail.com"},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "",
			ExpectedError:   "user already exists",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mockAuthService.EXPECT().SignupUser(c.Request).AnyTimes().Return(newResp, nil)
			testServer.handler.AuthService = mockAuthService
			data, err := json.Marshal(c.Request)
			require.NoError(t, err)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(data))
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			testServer.router.ServeHTTP(recorder, req)
			assert.Equal(t, recorder.Code, c.ExpectedCode)
			assert.Contains(t, recorder.Body.String(), c.ExpectedMessage)

		})
	}
}
