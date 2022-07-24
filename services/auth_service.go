package services

import (
	"errors"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/dto"
	apiError "github.com/decagonhq/meddle-api/errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const AccessTokenValidity = time.Minute * 20
const RefreshTokenValidity = time.Hour * 24

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService

// AuthService interface
type AuthService interface {
	LoginUser(request *dto.LoginRequest, secret string) (*dto.LoginResponse, *apiError.Error)
}

// authService struct
type authService struct {
	authRepo db.AuthRepository
}

// NewAuthService instantiate an authService
func NewAuthService(authRepo db.AuthRepository) AuthService {
	return &authService{
		authRepo,
	}
}

func (a *authService) LoginUser(loginRequest *dto.LoginRequest, secret string) (*dto.LoginResponse, *apiError.Error) {
	foundUser, err := a.authRepo.FindUserByEmail(loginRequest.Email)
	if err != nil {
		return nil, apiError.ErrNotFound
	}

	if err := foundUser.VerifyPassword(loginRequest.Password); err != nil {
		return nil, apiError.ErrInvalidPassword
	}

	accessToken, err := GenerateToken(foundUser.Email, secret)
	if err != nil {
		return nil, apiError.ErrInternalServerError
	}

	userResponse := dto.UserResponse{
		ID:          foundUser.ID,
		Name:        foundUser.Name,
		PhoneNumber: foundUser.PhoneNumber,
		Email:       foundUser.Email,
	}

	return foundUser.LoginUserToDto(userResponse, accessToken), nil
}

// GenerateToken generates only an access token
func GenerateToken(email string, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("empty secret")
	}
	// Generate claims
	claims := GenerateClaims(email)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateClaims(email string) jwt.MapClaims {

	accessClaims := jwt.MapClaims{
		"email":   email,
		"expired": time.Now().Add(AccessTokenValidity).Unix(),
	}

	return accessClaims
}
