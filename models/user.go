package models

import (
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserResponse
	AccessToken string
}

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
