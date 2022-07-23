package models

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type User struct {
	Model
	Name           string `json:"name" gorm:"not null" binding:"required"`
	Email          string `json:"email" gorm:"unique;not null" binding:"required,email"`
	PhoneNumber    string `json:"phone_number" binding:"required,e164" gorm:"not null" gorm:"unique"`
	Password       string `json:"password" binding:"required" gorm:"not null"`
	HashedPassword string `json:"-" gorm:"password"`
	IsEmailActive  bool   `json:"-"`
}

func (u *User) Validate() *errors.Error {
	err := validator.New().Struct(u)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)

		var s []string
		for _, v := range validationErrors {
			s = append(s, errors.NewFieldError(v).String())
		}
		return &errors.Error{
			Message: strings.Join(s, ","),
			Status:  http.StatusBadRequest,
		}
	}

	return nil
}
