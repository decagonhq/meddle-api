package services

import (
	"fmt"
	apiError "github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/jwt"
	"log"
	"net/http"
)

func (a *authService) SendEmailForPasswordReset(user *models.ForgotPassword) *apiError.Error {

	foundUser, err := a.authRepo.FindUserByEmail(user.Email)
	if err != nil {
		return apiError.New("email does not exist", http.StatusBadRequest)
	}
	token, err := GenerateToken(foundUser.Email, a.Config.JWTSecret)
	if err != nil {
		return apiError.New("", http.StatusInternalServerError)
	}
	link := fmt.Sprintf("%s:%d/api/v1/password/reset/%s", a.Config.Host, a.Config.Port, token)
	log.Println("link: ", link)
	log.Println("token: ", token)
	body := "Please Click the link below to reset your password"
	title := "Password Reset Link"
	value := map[string]interface{}{}
	value["link"] = link
	err = a.mail.SendMail(user.Email, title, body, "", value)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return apiError.New("mail couldn't be sent", http.StatusServiceUnavailable)
	}
	return nil
}

func (a *authService) ResetPassword(reset *models.ResetPassword, token string) *apiError.Error {
	err := models.ValidatePassword(reset.Password)
	if err != nil {
		return apiError.New("", http.StatusBadRequest)
	}
	if reset.Password != reset.ConfirmPassword {
		return apiError.New("password does not match", http.StatusBadRequest)
	}
	var user models.User
	user.Password = reset.Password
	user.HashedPassword, err = GenerateHashPassword(user.Password)
	if err != nil {
		return apiError.New("", http.StatusInternalServerError)
	}
	err = a.authRepo.IsTokenInBlacklist(token)
	if err != nil {
		return apiError.New("expired link", http.StatusUnauthorized)
	}
	claims, err := jwt.ValidateAndGetClaims(token, a.Config.JWTSecret)
	if err != nil {
		return apiError.New("invalid link", http.StatusUnauthorized)
	}
	email := claims["email"].(string)
	errr := a.authRepo.UpdatePassword(user.HashedPassword, email)
	if errr != nil {
		return apiError.New("", http.StatusInternalServerError)
	}
	accBlacklist := &models.BlackList{
		Email: email,
		Token: token,
	}
	if err := a.authRepo.AddToBlackList(accBlacklist); err != nil {
		return apiError.New("", http.StatusInternalServerError)
	}
	return nil
}
