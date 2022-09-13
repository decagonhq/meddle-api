package server

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *Server) authorizeNotificationsForDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenArgument models.AddNotificationTokenArgs

		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}

		userId := user.ID
		if err := decode(c, &tokenArgument); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		tokenArgument.UserID = userId
		deviceToken, err := s.PushNotification.AuthorizeNotification(&tokenArgument)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "device authorized to receive notification successfully", http.StatusCreated, deviceToken, nil)
		registrationTokens, err := s.PushNotification.GetSingleUserDeviceTokens(int(tokenArgument.UserID))
		if err != nil {
			err.Respond(c)
			return
		}
		message := fmt.Sprintf("welcome %v ,your device has been enbled", user.Name)
		pushPayload := &models.PushPayload{
			Title: "Welcome Message",
			Body:  message,
		}
		time.Sleep(time.Second * 3)
		_, err = s.PushNotification.SendPushNotification(registrationTokens, pushPayload)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "Your device is enabled", http.StatusOK, message, nil)
	}
}
