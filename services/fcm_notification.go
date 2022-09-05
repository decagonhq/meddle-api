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
	"google.golang.org/api/option"
	"log"
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services PushNotification

type PushNotifier interface {
	SendPushNotificationToSingleDevice(fcmClient *messaging.Client) (string, error)
	SendPushNotificationToMultipleDevice(fcmClient *messaging.Client) (*messaging.BatchResponse, error)
	FirebaseInit() error
	GetDecodedFireBaseKey() ([]byte, error)
	AuthorizeNotification(request *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, *errors.Error)
}

type notificationService struct {
	Conf             *config.Config
	notificationRepo db.NotificationRepository
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

func (m *notificationService) AuthorizeNotification(request *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, *errors.Error) {
	token, err := m.notificationRepo.AddNotificationToken(request)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return token, nil
}
