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

func (s *Server) handleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func GetClaimsFromToken(c *gin.Context) (string, *models.User, *errors.Error) {
	var tokenI, userI interface{}
	var tokenExists, userExists bool

	if tokenI, tokenExists = c.Get("access_token"); !tokenExists {
		// response.JSON(c, "", http.StatusForbidden, nil, errors.New("forbidden", http.StatusForbidden))
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}
	if userI, userExists = c.Get("user"); !userExists {
		// response.JSON(c, "", http.StatusForbidden, nil, errors.New("forbidden", http.StatusForbidden))
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
		token, user, err := GetClaimsFromToken(c)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}
		// TODO: check if token has not expired
		if token.exp < time.Now() {
			response.JSON(c, "successfully logged out", http.StatusOK, nil, nil)
			return
		} else {
			accBlacklist := &models.BlackList{
			Email: user.Email,
			Token: token,
		}
		if err := s.AuthRepository.AddToBlackList(accBlacklist); err != nil {
			log.Printf("can't add access token to blacklist: %v\n", err)
			response.JSON(c, "logout failed", http.StatusInternalServerError, nil, errors.New("can't add access token to blacklist", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "logout successful", http.StatusOK, nil, nil)
	}
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
