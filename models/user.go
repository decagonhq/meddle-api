package models

import (
	"errors"
	"github.com/decagonhq/meddle-api/dto"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Name           string `json:"name"`
	PhoneNumber    string `json:"phone_number"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	HashedPassword string `json:"-"`
	IsAgree        bool   `json:"is_agree"`
}

// VerifyPassword verifies the collected password with the user's hashed password
func (u *User) VerifyPassword(password string) error {
	if u.HashedPassword != "" && len(u.HashedPassword) == 0 {
		// Internal Server
		return errors.New("password is not set")
	}
	// Wrong Password
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

// LoginUserToDto responsible for creating a response object for the handleLogin handler
func (u *User) LoginUserToDto(user dto.UserResponse, token string) *dto.LoginResponse {
	return &dto.LoginResponse{
		UserResponse: user,
		AccessToken:  token,
	}
}
