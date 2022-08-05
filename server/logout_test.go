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

var mockRepository *mocks.MockAuthRepository

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	ctrl.Finish()
	mockRepository = mocks.NewMockAuthRepository(ctrl)

	return func() {
		defer ctrl.Finish()
	}
}

func Test_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := mocks.NewMockAuthService(ctrl)
	repo := mocks.NewMockAuthRepository(ctrl)

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
	conf.JWTSecret = "testSecret"
	token, err := services.GenerateToken(user.Email, conf.JWTSecret)

	s := &Server{
		Config:         conf,
		AuthRepository: repo,
		AuthService:    auth,
	}

	repo.EXPECT().AddToBlackList(&models.BlackList{Email: user.Email, Token: token}).Return(nil)
	repo.EXPECT().TokenInBlacklist(token).Return(false)
	repo.EXPECT().FindUserByEmail(user.Email).Return(user, nil)

	r := s.setupRouter()
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/logout", strings.NewReader(string(user.Email)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	r.ServeHTTP(resp, req)
	fmt.Println(resp.Body.String())
	assert.Equal(t, 200, resp.Code)
}
