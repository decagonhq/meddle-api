package server

import (
	"encoding/json"
	"github.com/decagonhq/meddle-api/db/mocks"
	"github.com/decagonhq/meddle-api/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	_ "gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuyerSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDb := mocks.NewMockAuthRepository(ctrl)
	h := &Server{AuthRepository: mockDb}
	r := h.setupRouter()

	user := models.User{
		Model:          models.Model{}, //CHECK
		Email:          "TOLU@gmail.com",
		HashedPassword: "$2a$10$/pUzdX5zckOpb1jhC1jJZ.mlfOCO4Xy5YKgUWPt8GwtlUFTaVtxeC",
		PhoneNumber:    "08166677888",
		Password:       "password",
		IsEmailActive:  true,
		IsAgree:        true,
	}

	mockDb.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)
	newUser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(string(newUser)))
	r.ServeHTTP(w, req)

	mockDb.EXPECT().FindUserByEmail(user.Email).Return(&user, nil)

	t.Run("Check if email exists", func(t *testing.T) {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(string(newUser)))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "email already exists")

	})

}
