package models

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/leebenson/conform"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Name           string `json:"name" binding:"required,min=2"`
	Email          string `json:"email" gorm:"unique;not null" binding:"required,email"`
	PhoneNumber    string `json:"phone_number" gorm:"unique" binding:"required,e164"`
	Password       string `json:"password,omitempty" gorm:"-" binding:"required,min=8,max=15"`
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

type UserResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type LogoutRequest struct {
	Email    string `json:"email" binding:"required,email"`
}
type LoginResponse struct {
	UserResponse
	AccessToken string
}
type LogoutResponse struct {
	UserResponse
	AccessToken string
}
// VerifyPassword verifies the collected password with the user's hashed password
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

// LoginUserToDto responsible for creating a response object for the handleLogin handler
func (u *User) LoginUserToDto(token string) *LoginResponse {
	return &LoginResponse{
		UserResponse: UserResponse{
			ID:          u.ID,
			Name:        u.Name,
			PhoneNumber: u.PhoneNumber,
			Email:       u.Email,
		},
		AccessToken: token,
	}
}
