package server

import (
	"bytes"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/decagonhq/meddle-api/db/mocks"
	"github.com/decagonhq/meddle-api/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignup(t *testing.T) {
	newReq := &dto.SignupRequest{
		Name:        "Tolu",
		PhoneNumber: "+2348163608141",
		Email:       "toluwase@gmail.com",
		Password:    "12345678",
	}
	newResp := &dto.SignupResponse{
		Name:        newReq.Name,
		PhoneNumber: newReq.PhoneNumber,
		Email:       newReq.Email,
		ID:          gofakeit.UintRange(1, 100),
	}
	cases := []struct {
		Name            string
		User            *dto.SignupRequest
		ExpectedCode    int
		ExpectedMessage string
		ExpectedError   string
		buildStubs      func(ctrl *mocks.MockAuthService)
		checkResponse   func(recorder *httptest.ResponseRecorder)
	}{
		{
			Name:            "Test Signup with correct details",
			User:            newReq,
			ExpectedCode:    http.StatusCreated,
			ExpectedMessage: "user created successfully",
			ExpectedError:   "",
		},
		{
			Name:            "Test Signup without no email",
			User:            &dto.SignupRequest{Name: "Tolu", PhoneNumber: "+2348163608141", Password: "12345678"},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "",
			ExpectedError:   "Email is invalid: ",
		},
		{
			Name:            "Test Signup with invalid fields",
			User:            &dto.SignupRequest{Name: "Tolu", PhoneNumber: "08141", Email: "tolut.a"},
			ExpectedCode:    http.StatusBadRequest,
			ExpectedMessage: "",
			ExpectedError:   "Email is invalid: toluwase.tt.com",
		},
		{
			Name:            "Test Signup with duplicate email address",
			User:            &dto.SignupRequest{Name: "Tolu", PhoneNumber: "08141", Email: "toluwase@gmail.com"},
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
			mockAuthService.EXPECT().SignupUser(c.User).AnyTimes().Return(newResp, nil)
			testServer.handler.AuthService = mockAuthService
			data, err := json.Marshal(c.User)
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
func TestServer_HandleSignup(t *testing.T) {
	signupRequest := &dto.SignupRequest{
		Name:        "Tolu",
		PhoneNumber: "+2348163608141",
		Email:       "toluwase@gmail.com",
		Password:    "12345678",
	}
	signupResponse := &dto.SignupResponse{
		Name:        signupRequest.Name,
		PhoneNumber: signupRequest.PhoneNumber,
		Email:       signupRequest.Email,
		ID:          gofakeit.UintRange(1, 100),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuthService := mocks.NewMockAuthService(ctrl)
	mockAuthService.EXPECT().SignupUser(signupRequest).AnyTimes().Return(signupResponse, nil)
	testServer.handler.AuthService = mockAuthService
	data, err := json.Marshal(signupRequest)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(data))
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	testServer.router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

}

//func TestSignup(t *testing.T) {
//	emailAddress := gofakeit.Email()
//	cases := []struct {
//		Name            string
//		User            *dto.SignupRequest
//		ExpectedCode    int
//		ExpectedMessage string
//		ExpectedError   string
//		buildStubs      func(store *mocks.MockAuthRepository)
//		checkResponse   func(recoder *httptest.ResponseRecorder)
//	}{
//		{
//			Name: "Test Signup with correct details",
//			User: &dto.SignupRequest{
//				Name:        gofakeit.FirstName() + " " + gofakeit.LastName(),
//				Password:    gofakeit.Password(true, true, true, true, false, 8),
//				Email:       emailAddress,
//				PhoneNumber: gofakeit.Phone(),
//			},
//			ExpectedCode:    http.StatusCreated,
//			ExpectedMessage: "signup successful",
//			ExpectedError:   "",
//		},
//		{
//			Name: "Test Signup without an email address",
//			User: &dto.SignupRequest{
//				Name:     gofakeit.FirstName() + " " + gofakeit.LastName(),
//				Password: gofakeit.Password(true, true, true, true, false, 8),
//			},
//			ExpectedCode:    http.StatusBadRequest,
//			ExpectedMessage: "",
//			ExpectedError:   "Email is invalid: ",
//		},
//		{
//			Name: "Test Signup with invalid fields",
//			User: &dto.SignupRequest{
//				Name:     gofakeit.FirstName() + " " + gofakeit.LastName(),
//				Password: gofakeit.Password(true, true, true, true, false, 8),
//				Email:    "tolu.tee.gmail.com",
//			},
//			ExpectedCode:    http.StatusBadRequest,
//			ExpectedMessage: "",
//			ExpectedError:   "Email is invalid: odohi.davidgmail.com",
//		},
//		{
//			Name: "Test Signup with duplicate email address",
//			User: &dto.SignupRequest{
//				Name:     gofakeit.FirstName() + " " + gofakeit.LastName(),
//				Password: gofakeit.Password(true, true, true, true, false, 8),
//				Email:    emailAddress,
//			},
//			ExpectedCode:    http.StatusBadRequest,
//			ExpectedMessage: "",
//			ExpectedError:   "user already exists",
//		},
//	}
//
//	for i := range cases {
//		tc := cases[i]
//
//		t.Run(tc.Name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			store := mocks.NewMockAuthRepository(ctrl)
//			//store := mockdb.NewMockStore(ctrl)
//			tc.buildStubs(store)
//
//			server := newTestServer(t, store)
//			recorder := httptest.NewRecorder()
//
//			// Marshal body data to JSON
//			data, err := json.Marshal(tc.body)
//			require.NoError(t, err)
//
//			url := "/accounts"
//			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
//			require.NoError(t, err)
//
//			tc.setupAuth(t, request, server.tokenMaker)
//			server.router.ServeHTTP(recorder, request)
//			tc.checkResponse(recorder)
//		})
//	}
//
//	for i, c := range cases {
//		t.Run(c.Name, func(t *testing.T) {
//			code, res, err := newRequest(http.MethodPost, "/api/v1/auth/signup", "", c.User)
//			require.NoError(t, err)
//
//			assert.Equal(t, c.ExpectedCode, code, "Expected code to be %d, got %d", c.ExpectedCode, code)
//			assert.Equal(t, c.ExpectedMessage, res["message"])
//			assert.Equal(t, c.ExpectedError, res["errors"])
//		})
//	}
//}

//func TestHandleSignup(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	mockDb := mocks.NewMockAuthRepository(ctrl)
//	h := &Server{AuthRepository: mockDb}
//	r := h.setupRouter()
//
//	user := models.User{
//		Email:         "TOLU@gmail.com",
//		PhoneNumber:   "08166677888",
//		Password:      "password",
//		IsEmailActive: false,
//	}
//
//	newUser, err := json.Marshal(user)
//	if err != nil {
//		t.Fail()
//	}
//
//	t.Run("Check if email or phone exists", func(t *testing.T) {
//		mockDb.EXPECT().FindUserByEmailOrPhoneNumber(user.Email, user.PhoneNumber).Return(&user, nil).AnyTimes()
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(string(newUser)))
//		r.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusNotFound, w.Code)
//		assert.Contains(t, w.Body.String(), "email or phone already exists")
//
//	})
//}
