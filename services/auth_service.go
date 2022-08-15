package services

import (
	"crypto/rand"
	"encoding/json"
	"encoding/json"
	"fmt"
	"github.com/decagonhq/meddle-api/services/jwt"
	"io/ioutil"
	"github.com/decagonhq/meddle-api/config"
	"io/ioutil"
	"math/rand"
	"net/http"

	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/decagonhq/meddle-api/db"
	apiError "github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	//"github.com/golang-jwt/jwt"
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService

// AuthService interface
type AuthService interface {
	LoginUser(request *models.LoginRequest) (*models.LoginResponse, *apiError.Error)
	SignupUser(request *models.User) (*models.User, *apiError.Error)
	FacebookSignInUser(token string) (*string, *apiError.Error)
	VerifyEmail(token string) error
	SendEmailForPasswordReset(user *models.ForgotPassword) *apiError.Error
	ResetPassword(user *models.ResetPassword, token string) *apiError.Error
	GoogleSignInUser(token string) (*string, *apiError.Error)
}

// authService struct
type authService struct {
	Config   *config.Config
	authRepo db.AuthRepository
	mail     Mailer
}

// NewAuthService instantiate an authService
func NewAuthService(authRepo db.AuthRepository, conf *config.Config, mailer Mailer) AuthService {
	return &authService{
		Config:   conf,
		authRepo: authRepo,
		mail:     mailer,
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
	token, err := jwt.GenerateToken(user.Email, a.Config.JWTSecret)
	if err != nil {
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	link := fmt.Sprintf("http://localhost:8080/verifyEmail/%s", token)
	subject := "Verify your email"
	body := "Please Click the link below to verify your email"
	templateName := "verifyEmail"
	err = a.mail.SendMail(user.Email, subject, body, templateName, map[string]interface{}{link: link})
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return nil, apiError.New("mail couldn't be sent", http.StatusServiceUnavailable)
	}
	return user, nil
}

func (a *authService) GoogleSignInUser(token string) (*string, *apiError.Error) {

	googleUserDetails, googleUserDetailsError := GetUserInfoFromGoogle(token)

	if googleUserDetailsError != nil {
		return nil, apiError.New(fmt.Sprintf("unable to get user details from google: %v", googleUserDetailsError), http.StatusUnauthorized)
	}

	authToken, authTokenError := a.GetSignInToken(googleUserDetails)

	if authTokenError != nil {
		return nil, apiError.New(fmt.Sprintf("unable sign in user: %v", authTokenError), http.StatusUnauthorized)
	}
	return &authToken, nil
}

// GetUserInfoFromGoogle will return information of user which is fetched from Google
func GetUserInfoFromGoogle(token string) (*models.GoogleUser, error) {
	var googleUserDetails *models.GoogleUser

	url := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token
	googleUserDetailsRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Google: %+v", err)
	}

	googleUserDetailsResponse, googleDetailsResponseError := http.DefaultClient.Do(googleUserDetailsRequest)
	if googleDetailsResponseError != nil {
		return nil, fmt.Errorf("error occurred while getting information from Google: %+v", googleDetailsResponseError)
	}

	body, err := ioutil.ReadAll(googleUserDetailsResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Google: %+v", err)
	}
	defer googleUserDetailsResponse.Body.Close()

	err = json.Unmarshal(body, &googleUserDetails)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Google: %+v", err)
	}

	return googleUserDetails, nil
}

// GetSignInToken Used for Signing In the Users
func (a *authService) GetSignInToken(googleUserDetails *models.GoogleUser) (string, error) {
	var result *models.User

	if googleUserDetails == nil {
		return "", fmt.Errorf("error: Google user details can't be empty")
	}

	if googleUserDetails.Email == "" {
		return "", fmt.Errorf("error: email can't be empty")
	}

	if googleUserDetails.Name == "" {
		return "", fmt.Errorf("error: google name can't be empty")
	}

	result, err := a.authRepo.FindUserByEmail(googleUserDetails.Email)
	if err != nil {
		return "", fmt.Errorf("error finding user: %+v", err)
	}

	if result == nil {
		result.Email = googleUserDetails.Email
		result.Name = googleUserDetails.Name
		_, err = a.authRepo.CreateUser(result)
		if err != nil {
			return "", fmt.Errorf("error occurred creating user: %+v", err)
		}
	}

	tokenString, err := GenerateToken(googleUserDetails.Email, a.Config.JWTSecret)

	if tokenString == "" {
		return "", fmt.Errorf("unable to generate Auth token: %+v", err)
	}

	return tokenString, nil
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

	accessToken, err := jwt.GenerateToken(foundUser.Email, a.Config.JWTSecret)
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

func GenerateRandomString() (string, error) {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%X", b)
	return s, nil
}

func (a *authService) VerifyEmail(token string) error {
	//validate token here
	err := a.authRepo.VerifyEmail(token)
	return err
}

func (a *authService) FacebookSignInUser(token string) (*string, *apiError.Error) {
	// rename function
	fbUserDetails, fbUserDetailsError := GetUserInfoFromFacebook(token)

	if fbUserDetailsError != nil {
		return nil, apiError.New(fmt.Sprintf("unable to get user details from facebook: %v", fbUserDetailsError), http.StatusUnauthorized)
	}

	authToken, authTokenError := a.GetSignInToken(fbUserDetails)

	if authTokenError != nil {
		return nil, apiError.New(fmt.Sprintf("unable sign in user: %v", authTokenError), http.StatusUnauthorized)
	}
	return &authToken, nil
}

// GetUserInfoFromFacebook will return information of user which is fetched from facebook
func GetUserInfoFromFacebook(token string) (*models.FacebookUser, error) {
	var fbUserDetails *models.FacebookUser
	facebookUserDetailsRequest, _ := http.NewRequest("GET", "https://graph.facebook.com/me?fields=name,email&access_token="+token, nil)
	facebookUserDetailsResponse, facebookUserDetailsResponseError := http.DefaultClient.Do(facebookUserDetailsRequest)

	if facebookUserDetailsResponseError != nil {
		return nil, fmt.Errorf("error occurred while getting information from Facebook: %+v", facebookUserDetailsResponseError)
	}
	body, err := ioutil.ReadAll(facebookUserDetailsResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Facebook: %+v", err)
	}
	defer facebookUserDetailsResponse.Body.Close()
	err = json.Unmarshal(body, &fbUserDetails)

	if err != nil {
		return nil, fmt.Errorf("error occurred while getting information from Facebook: %+v", err)
	}

	return fbUserDetails, nil
}

// GetSignInToken Used for Signing In the Users
func (a *authService) GetSignInToken(facebookUserDetails *models.FacebookUser) (string, error) {
	var result *models.User

	if facebookUserDetails == nil {
		return "", fmt.Errorf("error: facebook user details can't be empty")
	}

	if facebookUserDetails.Email == "" {
		return "", fmt.Errorf("error: email can't be empty")
	}

	if facebookUserDetails.Name == "" {
		return "", fmt.Errorf("error: name can't be empty")
	}

	result, err := a.authRepo.FindUserByEmail(facebookUserDetails.Email)
	if err != nil {
		return "", fmt.Errorf("error finding user: %+v", err)
	}

	if result == nil {
		result.Email = facebookUserDetails.Email
		result.Name = facebookUserDetails.Name
		_, err = a.authRepo.CreateUser(result)
		if err != nil {
			return "", fmt.Errorf("error occurred creating user: %+v", err)
		}
	}

	tokenString, err := jwt.GenerateToken(facebookUserDetails.Email, a.Config.JWTSecret)

	if tokenString == "" {
		return "", fmt.Errorf("unable to generate Auth token: %+v", err)
	}

	return tokenString, nil
}

func GenerateRandomString() (string, error) {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%X", b)
	return s, nil
}

