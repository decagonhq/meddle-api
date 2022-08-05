package server

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/decagonhq/meddle-api/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (s *Server) ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var foundUser models.ForgotPassword
		if err := c.ShouldBindJSON(&foundUser); err != nil {
			response.JSON(c, "error unmarshalling body", http.StatusBadRequest, nil, err)
			return
		}
		err := s.AuthService.ForgotPassword(&foundUser)
		if err != nil {
			response.JSON(c, "email was not sent", http.StatusBadRequest, nil, err)
			return
		}
		response.JSON(c, "link to reset password successfully sent", http.StatusCreated, nil, nil)
	}
}

func (s *Server) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pw models.ResetPassword
		if err := c.ShouldBindJSON(&pw); err != nil {
			response.JSON(c, "error unmarshalling body", http.StatusBadRequest, nil, err)
			return
		}
		err := models.ValidatePassword(pw.Password)
		if err != nil {
			response.JSON(c, "", errors.ErrBadRequest.Status, nil, err)
			return
		}
		if pw.Password != pw.ConfirmPassword {
			response.JSON(c, "password does not match", errors.ErrBadRequest.Status, nil, err)
		}
		var user models.User
		user.Password = pw.Password
		user.HashedPassword, err = services.GenerateHashPassword(user.Password)
		if err != nil {
			log.Printf("error generating password hash: %v", err.Error())
			response.JSON(c, "internal server error", errors.ErrInternalServerError.Status, nil, err)
		}
		token := c.Param("token")
		claims, errr := getClaims(token, s.Config.JWTSecret)
		if errr != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errr)
			return
		}
		if isTokenExpired(claims) {
			response.JSON(c, "your link expired, cant update password", http.StatusInternalServerError, nil, errr)
			return
		}
		email := claims["email"].(string)
		errr = s.AuthRepository.UpdatePassword(user.HashedPassword, email)
		if errr != nil {
			response.JSON(c, "cant update password", http.StatusInternalServerError, nil, errr)
			return
		}
		convertClaims, _ := claims["exp"].(int64)
		if convertClaims < time.Now().Unix() {
			accBlacklist := &models.BlackList{
				Email: email,
				Token: token,
			}
			if err := s.AuthRepository.AddToBlackList(accBlacklist); err != nil {
				log.Printf("can't add access token to blacklist: %v\n", err)
				response.JSON(c, "reset failed", http.StatusInternalServerError, nil, errors.New("can't add access token to blacklist", http.StatusInternalServerError))
				return
			}
		}
		response.JSON(c, "Reset successful, Login with your new password to continue", http.StatusCreated, nil, nil)
	}
}
