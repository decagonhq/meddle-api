package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
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
		_, err = s.PushNotification.AuthorizeNotification(&tokenArgument)
		if err != nil {
			err.Respond(c)
			return
		}
		go func() {
			time.Sleep(time.Second * 3)
			message := "We will remind you to take your medications when it's due."
			pushPayload := &models.PushPayload{
				Title: fmt.Sprintf("Hello %s 👋", user.Name),
				Body:  message,
				Data: map[string]string{
					"medication_id": "23",
				},
				Category: models.WelcomeCategory,
			}
			// fmt.Printf("notification payload: %v\ntokenArgument: %+v", pushPayload, tokenArgument)
			_, err = s.PushNotification.SendPushNotification([]string{tokenArgument.Token}, pushPayload)
			if err != nil {
				log.Printf("error sending notification: %v", err)
			}
		}()
		response.JSON(c, "device authorized to receive notification successfully", http.StatusCreated, nil, nil)
	}
}
