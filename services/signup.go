package services

import (
	"fmt"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// AuthService interface
type AuthService interface {
	SignupUser(request *models.User) (*models.User, error)
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

func (a *authService) SignupUser(request *models.User) (*models.User, error) {
	var user = &models.User{}
	errs := models.ValidateStruct(request)
	if len(errs) > 0 {
		return nil, errors.New(fmt.Sprintf("%s", errs), http.StatusBadRequest)
	}
	exist, err := a.authRepo.IsEmailExist(request.Email)
	if exist {
		return nil, errors.New("email already exist", http.StatusBadRequest)
	}
	exist, err = a.authRepo.IsPhoneExist(request.PhoneNumber)
	if exist {
		return nil, errors.New("phone already exist", http.StatusBadRequest)
	}
	hashedPassword, err := GenerateHashPassword(request.Password)
	if err != nil {
		return nil, err
	}
	user.HashedPassword = string(hashedPassword)
	user.Password = ""
	newUser := models.User{
		Email:          request.Email,
		PhoneNumber:    request.PhoneNumber,
		Name:           request.Name,
		HashedPassword: string(hashedPassword),
		Password:       "",
		IsEmailActive:  false,
	}
	user, err = a.authRepo.CreateUser(&newUser)
	if err != nil {
		return nil, err
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

func GenerateHashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
