package server

import (
	"log"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
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

		claims, errr := getClaims(token, s.Config.JWTSecret)
		if errr != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errr)
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

func (s *Server) HandleVerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenBack, ok := c.Get("token")
		if !ok {
			response.JSON(c, "", http.StatusBadRequest, nil, errors.New("token not found", http.StatusBadRequest))
			return
		}
		userid := c.Param("userid")
		token,_,_:= GetValuesFromContext(c)
		if tokenBack == "" || userid == "" {
			log.Printf("token or userid is empty\n")
		}
		if token != tokenBack {
			log.Printf("token is not equal to tokenBack\n")
		} else {
			if err := s.AuthRepository.VerifyEmail(token, userid); err != nil {
				log.Printf("can't verify email: %v\n", err)
				response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("can't verify email", http.StatusInternalServerError))
				return
			}
			response.JSON(c, "email verified successfully", http.StatusOK, nil, nil)
		}
	}
}