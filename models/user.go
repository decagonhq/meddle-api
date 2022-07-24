package models

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/leebenson/conform"
)

type User struct {
	Model
	Name           string `json:"name" conform:"trim" validate:"required,min=2,max=15"`
	Email          string `json:"email" gorm:"unique;not null" validate:"required,email"`
	PhoneNumber    string `json:"phone_number" gorm:"unique" validate:"required,e164"`
	Password       string `json:"password" conform:"trim" binding:"required" validate:"required,min=8,max=15"`
	HashedPassword string `json:"-" gorm:"password"`
	IsEmailActive  bool   `json:"-"`
}

func ValidateStruct(req interface{}) []error {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	err := validateWhiteSpaces(req)
	errs := translateError(err, trans)
	err = validate.Struct(req)
	errs = translateError(err, trans)
	return errs
}

func validateWhiteSpaces(data interface{}) error {
	return conform.Strings(data)
}

func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans) + "; ")
		errs = append(errs, translatedErr)
	}
	return errs
}
