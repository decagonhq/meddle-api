package services

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/go-co-op/gocron"
	"google.golang.org/api/option"
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services PushNotification

type PushNotifier interface {
	AuthorizeNotification(request *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, *errors.Error)
	GetSingleUserDeviceTokens(userId int) ([]string, *errors.Error)
	CheckIfThereIsNextMedication()
	CheckIfThereWillBeNextMedicationInTheNext15Minutes()
	SendPushNotification(registrationTokens []string, payload *models.PushPayload) (*messaging.Message, *errors.Error)
	NotificationsCronJob()
	NotificationsCronJobFor15MinutesEarly()
}

type notificationService struct {
	Conf             *config.Config
	notificationRepo db.NotificationRepository
	Client           *messaging.Client
}

// NewFirebaseCloudMessaging instantiates an FCM service
func NewFirebaseCloudMessaging(notificationRepo db.NotificationRepository, conf *config.Config) (PushNotifier, error) {
	firebaseApp, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(conf.GoogleApplicationCredentials))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var fcm notificationService

	fcm.Client, err = firebaseApp.Messaging(context.Background())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &notificationService{
		notificationRepo: notificationRepo,
		Conf:             conf,
		Client:           fcm.Client,
	}, nil
}

func (fcm *notificationService) AuthorizeNotification(request *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, *errors.Error) {
	token, err := fcm.notificationRepo.AddNotificationToken(request)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return token, nil
}

func (fcm *notificationService) GetSingleUserDeviceTokens(userid int) ([]string, *errors.Error) {
	tokens, err := fcm.notificationRepo.GetSingleUserDeviceTokens(userid)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return tokens, nil
}

// CheckIfThereIsNextMedication cron job
// check all currently due medication in db
func (fcm *notificationService) CheckIfThereIsNextMedication() {
	medicationNotifications, err := fcm.notificationRepo.GetAllNextMedicationsToSendNotifications()
	if err != nil {
		log.Println("could not get medications from db", err)
		return
	}

	//check db for all the time of notifications
	for _, medicationNotification := range medicationNotifications {
		go func(m models.Medication) {
			userId := m.UserID
			deviceTokens, err := fcm.notificationRepo.GetSingleUserDeviceTokens(int(userId))
			if err != nil {
				log.Printf("error retrieving device notification tokens: %v\n", err)
				return
			}

			if len(deviceTokens) == 0 {
				log.Printf("empty token list: %v\n", err)
				return
			}
			nextDosageTime := m.NextDosageTime.Add(time.Hour).Format(time.Kitchen)
			notification, err := fcm.SendPushNotification(deviceTokens, &models.PushPayload{
				Body:  fmt.Sprintf("%s is due by %v", m.Name, nextDosageTime),
				Title: fmt.Sprintf("Time to take %s", m.Name),
				Data: map[string]string{
					"medication_id": fmt.Sprintf("%v", m.ID),
				},
				Category: models.NextMedicationCategory,
				// ClickAction: "/user/medication/id?=" + strconv.Itoa(int((m.ID)),
			})
			if err != nil {
				log.Println("error sending notification", err)
				return
			}

			log.Println("logging notifications", notification)
		}(medicationNotification)
	}
}

func (fcm *notificationService) CheckIfThereWillBeNextMedicationInTheNext15Minutes() {
	medicationNotifications, err := fcm.notificationRepo.GetAllNextMedicationsToSendNotifications()
	if err != nil {
		log.Println("could not get medications from db", err)
		return
	}

	//check db for all the time of notifications
	for _, medicationNotification := range medicationNotifications {
		go func(m models.Medication) {
			userId := m.UserID
			deviceTokens, err := fcm.notificationRepo.GetSingleUserDeviceTokens(int(userId))
			if err != nil {
				log.Printf("error retrieving device notification tokens: %v\n", err)
				return
			}

			if len(deviceTokens) == 0 {
				log.Printf("empty token list: %v\n", err)
				return
			}
			if m.NextDosageTime.Add(-time.Minute*15) == time.Now() {
				notification, err := fcm.SendPushNotification(deviceTokens, &models.PushPayload{
					Body:  fmt.Sprintf("%s will be due in 15 minutes", m.Name),
					Title: fmt.Sprintf("%s pre-notification reminder", m.Name),
					Data: map[string]string{
						"medication_id": fmt.Sprintf("%v", m.ID),
					},
					Category: models.NextMedicationCategory,
					// ClickAction: "/user/medication/id?=" + strconv.Itoa(int((m.ID)),
				})
				if err != nil {
					log.Println("error sending notification", err)
					return
				}
				log.Println("logging notifications", notification)
			}

		}(medicationNotification)
	}
}

func (fcm *notificationService) SendPushNotification(registrationTokens []string, payload *models.PushPayload) (*messaging.Message, *errors.Error) {
	message := &messaging.Message{
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Category: string(payload.Category),
					Alert: &messaging.ApsAlert{
						Title: payload.Title,
						Body:  payload.Body,
					},
					Sound:            "default",
					ContentAvailable: true,
				},
			},
			FCMOptions: nil,
		},
		Data: payload.Data,
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: "https://imgur.com/a/hmt6Mx2",
		},

		Token: registrationTokens[0],
	}
	for _, t := range registrationTokens {
		log.Printf("tokens: %v", t)
	}
	res, err := fcm.Client.Send(context.Background(), message)
	if err != nil {
		log.Printf("error sending message: %v\n", err)
		return nil, errors.ErrInternalServerError
	}

	log.Printf("result from message 2: %v\n", res)
	// d, err := fcm.Client.SendMulticast(context.Background(), notification)
	// if err != nil {
	// 	return &messaging.MulticastMessage{}, errors.ErrInternalServerError
	// }

	// log.Printf("failure: %v, success: %v", d.FailureCount, d.SuccessCount)
	// for _, v := range d.Responses {
	// 	log.Printf("res: %+v", v)
	// }
	return message, nil
}

func (fcm *notificationService) NotificationsCronJob() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Minute().Do(func() {
		fcm.CheckIfThereIsNextMedication()
	})
	scheduler.StartBlocking()
}

func (fcm *notificationService) NotificationsCronJobFor15MinutesEarly() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Minute().Do(func() {
		fcm.CheckIfThereWillBeNextMedicationInTheNext15Minutes()
	})
	scheduler.StartBlocking()
}
