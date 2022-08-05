package services

import (
	"fmt"
	apiError "github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"log"
	"net/http"
)

func (a *authService) ForgotPassword(user *models.ForgotPassword) *apiError.Error {

	foundUser, err := a.authRepo.FindUserByEmail(user.Email)
	if err != nil {
		return apiError.New("email does not exist", http.StatusBadRequest)
	}
	token, err := GenerateToken(foundUser.Email, a.Config.JWTSecret)
	if err != nil {
		return apiError.New("could not generate link", http.StatusInternalServerError)
	}
	log.Println("TOKEN: ", token)

	_, err = a.mail.SendResetPassword(user.Email, fmt.Sprintf("http://localhost:8080/reset-password/%s", token))

	log.Println("LINK: ", fmt.Sprintf("http://localhost:8080/reset-password/%s", token))
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return apiError.New("mail couldn't be sent", http.StatusServiceUnavailable)
	}
	return apiError.New("", http.StatusCreated)
}
