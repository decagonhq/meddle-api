package services

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	apiError "github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"

	"github.com/golang-jwt/jwt"
)

const AccessTokenValidity = time.Minute * 20
const RefreshTokenValidity = time.Hour * 24

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService

// AuthService interface
type AuthService interface {
	LoginUser(request *models.LoginRequest) (*models.LoginResponse, *apiError.Error)
}

// authService struct
type authService struct {
	Config   *config.Config
	authRepo db.AuthRepository
}

// NewAuthService instantiate an authService
func NewAuthService(authRepo db.AuthRepository, conf *config.Config) AuthService {
	return &authService{
		Config:   conf,
		authRepo: authRepo,
	}
}

func (a *authService) LoginUser(loginRequest *models.LoginRequest) (*models.LoginResponse, *apiError.Error) {
	foundUser, err := a.authRepo.FindUserByEmail(loginRequest.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apiError.ErrNotFound
		} else {
			log.Printf("error from database: %v", err)
			return nil, apiError.ErrInternalServerError
		}
	}

	if err := foundUser.VerifyPassword(loginRequest.Password); err != nil {
		return nil, apiError.ErrInvalidPassword
	}

	accessToken, err := GenerateToken(foundUser.Email, a.Config.JWTSecret)
	if err != nil {
		log.Printf("error generating token %s", err)
		return nil, apiError.ErrInternalServerError
	}

	return foundUser.LoginUserToDto(accessToken), nil
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
		"email": email,
		"exp":   time.Now().Add(AccessTokenValidity).Unix(),
	}

	return accessClaims
}
