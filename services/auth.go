package services

import (
	"fmt"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/dto"
	"github.com/decagonhq/meddle-api/errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const AccessTokenValidity = time.Minute * 20
const RefreshTokenValidity = time.Hour * 24

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService

// AuthService interface
type AuthService interface {
	LoginUser(request *dto.LoginRequest, secret string) (*dto.LoginResponse, *errors.Error)
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

func (a *authService) LoginUser(loginRequest *dto.LoginRequest, secret string) (*dto.LoginResponse, *errors.Error) {
	foundUser, err := a.authRepo.FindUserByEmail(loginRequest.Email)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	if err := foundUser.VerifyPassword(loginRequest.Password); err != nil {
		//errors.New()
		return nil, errors.ErrUnauthorized
	}

	accessToken, err := GenerateToken(foundUser.Email, &secret)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	userResponse := dto.UserResponse{
		ID:          foundUser.ID,
		Name:        foundUser.Name,
		PhoneNumber: foundUser.PhoneNumber,
		Email:       foundUser.Email,
	}

	return foundUser.LoginUserToDto(userResponse, *accessToken), nil
}

// GetTokenFromHeader returns the token string in the authorization header
func GetTokenFromHeader(c *gin.Context) string {
	authHeader := c.Request.Header.Get("Authorization")
	if len(authHeader) > 8 {
		return authHeader[7:]
	}
	return ""
}

// verifyAccessToken verifies a token
func verifyToken(tokenString *string, claims jwt.MapClaims, secret *string) (*jwt.Token, error) {
	parser := &jwt.Parser{SkipClaimsValidation: true}
	return parser.ParseWithClaims(*tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(*secret), nil
	})
}

// AuthorizeToken check if a refresh token is valid
func AuthorizeToken(token *string, secret *string) (*jwt.Token, jwt.MapClaims, error) {
	if token != nil && *token != "" && secret != nil && *secret != "" {
		claims := jwt.MapClaims{}
		token, err := verifyToken(token, claims, secret)
		if err != nil {
			return nil, nil, err
		}
		return token, claims, nil
	}
	return nil, nil, fmt.Errorf("empty token or secret")
}

// GenerateToken generates only an access token
func GenerateToken(email string, secret *string) (*string, error) {
	// Generate claims
	claims := GenerateClaims(email)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(*secret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func GenerateClaims(email string) jwt.MapClaims {

	accessClaims := jwt.MapClaims{
		"email":   email,
		"expired": time.Now().Add(AccessTokenValidity).Unix(),
	}

	return accessClaims
}
