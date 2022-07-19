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

func TestHandleSignup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDb := mocks.NewMockAuthRepository(ctrl)
	h := &Server{AuthRepository: mockDb}
	r := h.setupRouter()

	user := models.User{
		Model:         models.Model{},
		Email:         "TOLU@gmail.com",
		PhoneNumber:   "08166677888",
		Password:      "password",
		IsEmailActive: false,
	}

	newUser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
	}

	t.Run("Check if email or phone exists", func(t *testing.T) {
		mockDb.EXPECT().FindUserByEmailOrPhoneNumber(user.Email, user.PhoneNumber).Return(&user, nil).AnyTimes()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", strings.NewReader(string(newUser)))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "email or phone already exists")

	})
}
