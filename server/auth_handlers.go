package server

import (
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/services"
	"log"
	"net/http"

	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		HashedPassword, err := services.GenerateHashPassword(user.Password)
		user.HashedPassword = string(HashedPassword)
		if err != nil {
			log.Printf("hash password err: %v\n", err)
			response.JSON(c, "hashed password was not be generated successfully", http.StatusInternalServerError, nil, nil)
			return
		}
		_, err = s.AuthRepository.FindUserByEmail(user.Email)
		if err == nil {
			response.JSON(c, "email already exists", http.StatusNotFound, err, nil)
			return
		}
		_, err = s.AuthRepository.FindUserByPhoneNumber(user.PhoneNumber)
		if err == nil {
			response.JSON(c, "phone number already exists", http.StatusNotFound, err, nil)
			return
		}

		_, err = s.AuthRepository.FindUserByPhoneNumber(user.PhoneNumber)
		if err == nil {
			response.JSON(c, "phone number already exists", http.StatusNotFound, err, nil)
			return
		}

		_, err = s.AuthRepository.CreateUser(&user)
		if err != nil {
			log.Printf("create user err: %v\n", err)
			response.JSON(c, "", http.StatusInternalServerError, err, nil)
			return
		}
		response.JSON(c, "user created successfully", http.StatusOK, user, nil)
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
