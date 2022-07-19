package server

import (
	"github.com/decagonhq/meddle-api/dto"
	"github.com/decagonhq/meddle-api/errors"
	"gorm.io/gorm"
	"net/http"

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
		var loginRequest *dto.LoginRequest
		if err := c.ShouldBindJSON(loginRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}

		userResponse, err := s.AuthService.LoginUser(loginRequest, s.Config.JWTSecret)
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				response.JSON(c, "user not found", http.StatusNotFound, nil, err)
				return
			case errors.InValidPasswordError:
				response.JSON(c, "invalid password", http.StatusUnauthorized, nil, err)
				return
			default:
				response.JSON(c, "internal server error", http.StatusInternalServerError, nil, err)
				return
			}
		}
		response.JSON(c, "successful", http.StatusOK, userResponse, nil)
	}
}

func (s *Server) handleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
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
