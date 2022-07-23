package server

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"net/http"

	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Server) HandleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			response.JSON(c, "bad request", http.StatusBadRequest, err.Error(), nil)
			return
		}
		userResponse, err := s.AuthService.SignupUser(&user)
		if err != nil {
			response.JSON(c, err.Error(), http.StatusBadRequest, nil, &errors.Error{Message: err.Error(), Status: http.StatusBadRequest})
			return
		}
		response.JSON(c, "user created successfully", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
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
