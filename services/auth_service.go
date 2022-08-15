package services

import (
	"fmt"
	"net/http"

	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	apiError "github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/jwt"
	jwToken "github.com/golang-jwt/jwt"
)

const AccessTokenValidity = time.Hour * 24
const RefreshTokenValidity = time.Hour * 24

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService

// AuthService interface
type AuthService interface {
	LoginUser(request *models.LoginRequest) (*models.LoginResponse, *apiError.Error)
	SignupUser(request *models.User) (*models.User, *apiError.Error)
	VerifyEmail(token string) error
	SendEmailForPasswordReset(user *models.ForgotPassword) *apiError.Error
	ResetPassword(user *models.ResetPassword, token string) *apiError.Error
}

// authService struct
type authService struct {
	Config   *config.Config
	authRepo db.AuthRepository
	mail    Mailer
}

// NewAuthService instantiate an authService
func NewAuthService(authRepo db.AuthRepository, conf *config.Config, mailer Mailer) AuthService {
	return &authService{
		Config:   conf,
		authRepo: authRepo,
		mail:  mailer,
	}
}

func (a *authService) SignupUser(user *models.User) (*models.User, *apiError.Error) {
	err := a.authRepo.IsEmailExist(user.Email)

	if err != nil {
		// FIXME: return the proper error message from the function
		// TODO: handle internal server error later
		return nil, apiError.New("email already exist", http.StatusBadRequest)
	}
	err = a.authRepo.IsPhoneExist(user.PhoneNumber)
	if err != nil {
		return nil, apiError.New("phone already exist", http.StatusBadRequest)
	}
	user.HashedPassword, err = GenerateHashPassword(user.Password)

	if err != nil {
		log.Printf("error generating password hash: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	user.Password = ""
	user.IsEmailActive = false
	user, err = a.authRepo.CreateUser(user)

	if err != nil {
		log.Printf("unable to create user: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	token, err := GenerateToken(user.Email, a.Config.JWTSecret)
	if err != nil {
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}

	link := fmt.Sprintf("%s:%d/api/v1/verifyEmail/%s", a.Config.Host, a.Config.Port, token)
	value := map[string]interface{}{}
	value["link"] = link
	subject := "Verify your email"
	body := "Please Click the link below to verify your email"
	templateName := "verifyEmail"
	err = a.mail.SendMail(user.Email,subject, body,templateName, value)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return nil, apiError.New("mail couldn't be sent", http.StatusServiceUnavailable)
	}

	return user, nil
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
func verifyToken(tokenString *string, claims jwToken.MapClaims, secret *string) (*jwToken.Token, error) {
	parser := &jwToken.Parser{SkipClaimsValidation: true}
	return parser.ParseWithClaims(*tokenString, claims, func(token *jwToken.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwToken.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(*secret), nil
	})
}

// AuthorizeToken check if a refresh token is valid
func AuthorizeToken(token *string, secret *string) (*jwToken.Token, jwToken.MapClaims, error) {
	if token != nil && *token != "" && secret != nil && *secret != "" {
		claims := jwToken.MapClaims{}
		token, err := verifyToken(token, claims, secret)
		if err != nil {
			return nil, nil, err
		}
		return token, claims, nil
	}
	return nil, nil, fmt.Errorf("empty token or secret")
}

func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
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
	token := jwToken.NewWithClaims(jwToken.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateClaims(email string) jwToken.MapClaims {
	accessClaims := jwToken.MapClaims{
		"email": email,
		"exp":   time.Now().Add(AccessTokenValidity).Unix(),
	}
	return accessClaims
}

func (a *authService) VerifyEmail(token string) error {
	claims, err := jwt.ValidateAndGetClaims(token, a.Config.JWTSecret)
	if err != nil {
		return apiError.New("invalid link", http.StatusUnauthorized)
	}
	email := claims["email"].(string)
	err = a.authRepo.VerifyEmail(email, token)
	return err
}