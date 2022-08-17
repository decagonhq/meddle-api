package server

import (
	"net/http"

	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Server) SendEmailForPasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		var foundUser models.ForgotPassword
		if err := decode(c, &foundUser); err != nil {
			response.JSON(c, "error unmarshalling body", http.StatusBadRequest, nil, err)
			return
		}
		err := s.AuthService.SendEmailForPasswordReset(&foundUser)
		if err != nil {
			response.JSON(c, "email was not sent", http.StatusBadRequest, nil, err)
			return
		}
		response.JSON(c, "link to reset password successfully sent", http.StatusOK, nil, nil)
	}
}

func (s *Server) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var password models.ResetPassword
		if err := decode(c, &password); err != nil {
			response.JSON(c, "error unmarshalling body", http.StatusBadRequest, nil, err)
			return
		}
		err := s.AuthService.ResetPassword(&password, c.Param("token"))
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "password reset successfully, Login to continue", http.StatusCreated, nil, nil)
	}
}
