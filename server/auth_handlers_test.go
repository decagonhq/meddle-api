package server

import (
	"github.com/decagonhq/meddle-api/models"
	_ "gorm.io/gorm"
	"testing"
)

func TestBuyerSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDb := mock_database.NewMockDB(ctrl)
	h := &handlers.Handler{DB: mockDb}
	r, _ := router.SetupRouter(h)

	user := models.User{
		Model:          models.Model{}, //CHECK
		Email:          "garber@gmail.com",
		HashedPassword: "", //CHECK
		PhoneNumber:    "08022334455",
		Password:       "password",
		IsAgree:        true,
	}

}
