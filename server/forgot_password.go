package server

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/decagonhq/meddle-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (s *Server) SendEmailForPasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		var foundUser models.ForgotPassword
		if err := c.ShouldBindJSON(&foundUser); err != nil {
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
		if err := c.ShouldBindJSON(&password); err != nil {
			response.JSON(c, "error unmarshalling body", http.StatusBadRequest, nil, err)
			return
		}
		err := models.ValidatePassword(password.Password)
		if err != nil {
			response.JSON(c, "", errors.ErrBadRequest.Status, nil, err)
			return
		}
		if password.Password != password.ConfirmPassword {
			response.JSON(c, "password does not match", errors.ErrBadRequest.Status, nil, err)
			return
		}
		var user models.User
		user.Password = password.Password
		user.HashedPassword, err = services.GenerateHashPassword(user.Password)
		if err != nil {
			log.Printf("error generating password hash: %v", err.Error())
			response.JSON(c, "internal server error", errors.ErrInternalServerError.Status, nil, err)
			return
		}
		token := c.Param("token")
		err = s.AuthRepository.IsTokenInBlacklist(token)
		if err != nil {
			response.JSON(c, "expired token, Please request a new password reset link", http.StatusUnauthorized, nil, nil)
			return
		}
		//TODO Refactor the test server, remove repository from the actual server
		//getClaims function contains verifyToken function
		//where token validity is verified
		claims, errr := getClaims(token, s.Config.JWTSecret)
		if errr != nil {
			response.JSON(c, "invalid link, please try again", http.StatusUnauthorized, nil, errr)
			return
		}
		err = claims.Valid()
		if err != nil {
			response.JSON(c, "your token has expired, cant update password, Request a new password reset link", http.StatusUnauthorized, nil, errr)
			return
		}
		email := claims["email"].(string)
		errr = s.AuthRepository.UpdatePassword(user.HashedPassword, email)
		if errr != nil {
			response.JSON(c, "An error occurred, try again", http.StatusInternalServerError, nil, errr)
			return
		}
		accBlacklist := &models.BlackList{
			Email: email,
			Token: token,
		}
		if err := s.AuthRepository.AddToBlackList(accBlacklist); err != nil {
			log.Printf("can't add access token to blacklist: %v\n", err)
			response.JSON(c, "request a new link", http.StatusInternalServerError, nil, errors.New("", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "Reset successful, Login with your new password to continue", http.StatusCreated, nil, nil)
	}
}
