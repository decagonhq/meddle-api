package server

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"log"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {

		if tokenI, exists := c.Get("access_token"); exists {
			if userI, exists := c.Get("user"); exists {
				if user, ok := userI.(*models.User); ok {
					if accessToken, ok := tokenI.(string); ok {
						accBlacklist := &models.BlackList{
							Email:     user.Email,
							Model: models.Model{CreatedAt: time.Now().Unix()},
							Token:     accessToken,
						}

						err := s.AuthRepository.AddToBlackList(accBlacklist)
						if err != nil {
							log.Printf("can't add access token to blacklist: %v\n", err)
							response.JSON(c, "logout failed", http.StatusInternalServerError, nil, errors.New("can't add access token to blacklist", http.StatusInternalServerError))
							return
						}
						response.JSON(c, "logout successful", http.StatusOK, nil, nil)
						return
					}
				}
			}
		}
		response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("can't get info from context", http.StatusInternalServerError))
		return
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenI , _:= c.Get("access_token")
		userI, _ := c.Get("user")
		 if tokenI == "access_token" && userI == "user" {
			 return
		 }
		user := userI.(*models.User)
		if accessToken, ok := tokenI.(string); ok {
			accBlacklist := &models.BlackList{
				Email:     user.Email,
				Token:     accessToken,
				Model:  models.Model{CreatedAt: time.Now().Unix()},
				}
				err := s.AuthRepository.AddToBlackList(accBlacklist)
				if err != nil {
					log.Printf("can't add access token to blacklist: %v\n", err)
					response.JSON(c, "logout failed", http.StatusInternalServerError, nil, errors.New("can't add access token to blacklist", http.StatusInternalServerError))
						return
					}
					response.JSON(c, "logout successful", http.StatusOK, nil, nil)
					return
			}
		response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("can't get info from context", http.StatusInternalServerError))
		return
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
