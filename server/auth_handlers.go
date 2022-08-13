package server

import (
	"context"
	"encoding/base64"
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/services"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/jwt"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Server) HandleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.AuthService.SignupUser(&user)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "user created successfully", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			response.JSON(c, "", errors.ErrBadRequest.Status, nil, err)
			return
		}
		userResponse, err := s.AuthService.LoginUser(&loginRequest)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}
		response.JSON(c, "login successful", http.StatusOK, userResponse, nil)
	}
}

func (s *Server) HandleGoogleOauthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetGoogleOAuthConfig(s.Config.GoogleClientID, s.Config.GoogleClientSecret, s.Config.GoogleRedirectURL)
		log.Println("conf", conf)
		log.Println("clientID", s.Config.GoogleClientID)
		log.Println("clientSecret", s.Config.GoogleClientSecret)
		log.Println("redirectURL", s.Config.GoogleRedirectURL)
		s.Config.OauthStateString, _ = services.GenerateRandomString()
		url := conf.AuthCodeURL(s.Config.OauthStateString, oauth2.AccessTypeOnline)
		log.Println("url: ", url)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (s *Server) HandleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		var state = c.Query("state")
		var code = c.Query("code")

		if state != s.Config.OauthStateString {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid login", http.StatusUnauthorized))
			return
		}

		var oauth2Config = config.GetGoogleOAuthConfig(s.Config.GoogleClientID, s.Config.GoogleClientID, s.Config.GoogleRedirectURL)

		token, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil || token == nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid token", http.StatusUnauthorized))
			return
		}

		authToken, errr := s.AuthService.GoogleSignInUser(token.AccessToken)
		if errr != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid authToken", http.StatusUnauthorized))
			return
		}

		response.JSON(c, "", http.StatusOK, authToken, nil)
	}
}

func GetValuesFromContext(c *gin.Context) (string, *models.User, *errors.Error) {
	var tokenI, userI interface{}
	var tokenExists, userExists bool

	if tokenI, tokenExists = c.Get("access_token"); !tokenExists {
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}
	if userI, userExists = c.Get("user"); !userExists {
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}

	token, ok := tokenI.(string)
	if !ok {
		return "", nil, errors.New("internal server error", http.StatusInternalServerError)
	}
	user, ok := userI.(*models.User)
	if !ok {
		return "", nil, errors.New("internal server error", http.StatusInternalServerError)
	}
	return token, user, nil
}

func (s *Server) handleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, user, err := GetValuesFromContext(c)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}
		claims, errr := jwt.ValidateAndGetClaims(token, s.Config.JWTSecret)
		if errr != nil {
			response.JSON(c, "", http.StatusUnauthorized, nil, errr)
			return
		}
		convertClaims, _ := claims["exp"].(int64) //jwt pkg to validate
		if convertClaims < time.Now().Unix() {
			accBlacklist := &models.BlackList{
				Email: user.Email,
				Token: token,
			}
			if err := s.AuthRepository.AddToBlackList(accBlacklist); err != nil {
				log.Printf("can't add access token to blacklist: %v\n", err)
				response.JSON(c, "logout failed", http.StatusInternalServerError, nil, errors.New("can't add access token to blacklist", http.StatusInternalServerError))
				return
			}
		}
		response.JSON(c, "logout successful", http.StatusOK, nil, nil)

	}
}

func (s *Server) handleGetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleUpdateUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleShowProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(2 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	return state
}

func (s *Server) HandleVerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		paramToken, ok := c.Get("token")
		if !ok {
			response.JSON(c, "", http.StatusBadRequest, nil, errors.New("token not found", http.StatusBadRequest))
			return
		}
		err := s.AuthService.VerifyEmail(paramToken.(string))
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
		}
	}
}