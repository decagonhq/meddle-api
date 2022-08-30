package services

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/decagonhq/meddle-api/config"
	"google.golang.org/api/option"
	"log"
)

type PushNotification interface {
	SendPushNotificationToSingleDevice(fcmClient *messaging.Client) (string, error)
	SendPushNotificationToMultipleDevice(fcmClient *messaging.Client) (*messaging.BatchResponse, error)
	FirebaseInit() error
	GetDecodedFireBaseKey() ([]byte, error)
}

type FirebaseCloudMessaging struct {
	Conf *config.Config
}

// NewFirebaseCloudMessaging instantiates an FCM service
func NewFirebaseCloudMessaging(conf *config.Config) PushNotification {
	return &FirebaseCloudMessaging{
		Conf: conf,
	}
}

func (fcm *FirebaseCloudMessaging) GetDecodedFireBaseKey() ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(fcm.Conf.FirebaseAuthKey)
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}

func (fcm *FirebaseCloudMessaging) FirebaseInit() error {
	decodedKey, err := fcm.GetDecodedFireBaseKey()
	if err != nil {
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
		return err
	}
	return nil
}

func (fcm *FirebaseCloudMessaging) SendPushNotificationToSingleDevice(fcmClient *messaging.Client) (string, error) {
	response, err := fcmClient.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Medication Alert!!!",
			Body:  "You have a new medication to take",
		},
		Token: "sample-device-token", // it's a single device token
	})

	if err != nil {
		return "", err
	}
	return response, nil
}

func (fcm *FirebaseCloudMessaging) SendPushNotificationToMultipleDevice(fcmClient *messaging.Client) (*messaging.BatchResponse, error) {
	response, err := fcmClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: "Medication Alert!!!",
			Body:  "You have a new medication to take",
		},
		Tokens: []string{}, // an array of device tokens should be passed here
	})

	if err != nil {
		return &messaging.BatchResponse{}, err
	}
	log.Println("Response success count : ", response.SuccessCount)
	log.Println("Response failure count : ", response.FailureCount)
	return response, nil
}
