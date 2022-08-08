package server

import (
	"fmt"
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//var mockRepository *mocks.MockAuthRepository

func Test_MedicationDetail(t *testing.T){
	ctrl := gomock.NewController(t)
	auth := mocks.NewMockAuthService(ctrl)
	repo := mocks.NewMockAuthRepository(ctrl)
	med := mocks.NewMockMedicationService(ctrl)


	conf, err := config.Load()
	if err != nil {
		t.Error(err)
	}
	user := &models.User{
		Name:        "Tolu",
		PhoneNumber: "+2348163608141",
		Email:       "toluwase@gmail.com",
		Password:    "12345678",
	}
	medication := &models.Medication{
		Duration:               3,
		MedicationPrescribedBy: "ken",
		UserID:                 1,
	}
	conf.JWTSecret = "testSecret"
	token, err := services.GenerateToken(user.Email, conf.JWTSecret)

	s := &Server{
		Config:         conf,
		AuthRepository: repo,
		AuthService:    auth,
		MedicationService: med,
	}

	//repo.EXPECT().AddToBlackList(&models.BlackList{Email: user.Email, Token: token}).Return(nil)
	repo.EXPECT().TokenInBlacklist(token).Return(false)
	med.EXPECT().GetMedicationDetail(uint(1), user.ID).Return(medication, nil)
	repo.EXPECT().FindUserByEmail(user.Email).Return(user, nil)


	r := s.setupRouter()
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/medications/1", strings.NewReader(""))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	r.ServeHTTP(resp, req)
	fmt.Println(resp.Body.String())
	assert.Equal(t, 200, resp.Code)
}
