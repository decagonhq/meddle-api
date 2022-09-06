package services

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/go-co-op/gocron"
	"google.golang.org/api/option"
	"log"
	"strconv"
	"time"
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services PushNotification

type PushNotifier interface {
	SendPushNotificationToSingleDevice(fcmClient *messaging.Client) (string, error)
	SendPushNotificationToMultipleDevice(fcmClient *messaging.Client) (*messaging.BatchResponse, error)
	FirebaseInit() error
	GetDecodedFireBaseKey() ([]byte, error)
	AuthorizeNotification(request *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, *errors.Error)
	CheckIfThereIsNextMedication()
	SendPushNotification(registrationTokens []string, payload *models.PushPayload) error
	NotificationsCronJob()
}

type notificationService struct {
	Conf             *config.Config
	notificationRepo db.NotificationRepository
	Client           *messaging.Client
}

// NewFirebaseCloudMessaging instantiates an FCM service
func NewFirebaseCloudMessaging(notificationRepo db.NotificationRepository, conf *config.Config) PushNotifier {
	return &notificationService{
		notificationRepo: notificationRepo,
		Conf:             conf,
	}
}

func (fcm *notificationService) GetDecodedFireBaseKey() ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(fcm.Conf.FirebaseAuthKey)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return decodedKey, nil
}

func (fcm *notificationService) FirebaseInit() error {
	decodedKey, err := fcm.GetDecodedFireBaseKey()
	if err != nil {
		log.Println(err)
		return err
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	// Initialize firebase app
	app, err := firebase.NewApp(context.Background(), nil, opts...)
	if err != nil {
		log.Printf("Error in initializing firebase app: %s\n", err)
		return err
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	_, err = fcm.SendPushNotificationToSingleDevice(fcmClient)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (fcm *notificationService) SendPushNotificationToSingleDevice(fcmClient *messaging.Client) (string, error) {
	response, err := fcmClient.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Medication Alert!!!",
			Body:  "You have a new medication to take",
		},
		Token: "sample-device-token", // it's a single device token
	})

	if err != nil {
		log.Println(err)
		return "", err
	}
	return response, nil
}

func (fcm *notificationService) SendPushNotificationToMultipleDevice(fcmClient *messaging.Client) (*messaging.BatchResponse, error) {
	response, err := fcmClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: "Medication Alert!!!",
			Body:  "You have a new medication to take",
		},
		Tokens: []string{}, // an array of device tokens should be passed here
	})

	if err != nil {
		log.Println(err)
		return &messaging.BatchResponse{}, err
	}
	log.Println("Response success count : ", response.SuccessCount)
	log.Println("Response failure count : ", response.FailureCount)
	return response, nil
}

func (fcm *notificationService) AuthorizeNotification(request *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, *errors.Error) {
	token, err := fcm.notificationRepo.AddNotificationToken(request)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return token, nil
}

// CheckIfThereIsNextMedication cron job
//check all currently due medication in db
func (fcm *notificationService) CheckIfThereIsNextMedication() {
	medicationNotifications, err := fcm.notificationRepo.GetAllNextMedicationsToSendNotifications()
	if err != nil {
		log.Println("could not get medications from db", err)
		return
	}

	//check db for all the time of notifications
	for i := 0; i < len(medicationNotifications); i++ {
		go func(i int) {
			deviceTokens, err := fcm.notificationRepo.GetNotificationTokens(int((medicationNotifications)[i].UserID))
			if err != nil {
				log.Printf("error retrieving device notification tokens: %v\n", err)
				return
			}

			if len(deviceTokens) == 0 {
				return
			}
			//fcm.SendPushNotificationToSingleDevice()
			err = fcm.SendPushNotification(deviceTokens, &models.PushPayload{
				Body:  "'" + (medicationNotifications)[i].Name + "' is due in___",
				Title: (medicationNotifications)[i].Name,
				Data: map[string]string{
					"link": "/user/medication/id?=" + strconv.Itoa(int((medicationNotifications)[i].ID)),
				},
				ClickAction: "/user/medication/id?=" + strconv.Itoa(int((medicationNotifications)[i].ID)),
			})
			if err != nil {
				log.Println("error sending notification")
				return
			}

		}(i)

		////send notification 10 minutes before time
		//if singleMedication.NextDosageTime == time.Now().Add(-time.Second*10) {
		//	//send notification based on the user
		//}
		////send notification at exact time
		//if singleMedication.NextDosageTime == time.Now() {
		//	//send notification based on the user NOW
		//}
	}
}

func (fcm *notificationService) SendPushNotification(registrationTokens []string, payload *models.PushPayload) error {
	notification := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: "https://i.imgur.com/mGOXXII.png",
		},
		Data: payload.Data,
		Webpush: &messaging.WebpushConfig{
			Data: payload.Data,
			Notification: &messaging.WebpushNotification{
				Title:   payload.Title,
				Body:    payload.Body,
				Icon:    "https://i.imgur.com/mGOXXII.png",
				Vibrate: []int{200, 100, 200},
			},
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Color:                 "#4C51BF",
				ClickAction:           payload.ClickAction,
				DefaultSound:          true,
				DefaultVibrateTimings: true,
				DefaultLightSettings:  true,
			},
			Data: payload.Data,
		},
		Tokens: registrationTokens,
	}

	_, err := fcm.Client.SendMulticast(context.Background(), notification)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func (fcm *notificationService) NotificationsCronJob() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Hour().Do(func() {
		fcm.CheckIfThereIsNextMedication()
	})
}
