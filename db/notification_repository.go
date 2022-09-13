package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../mocks/notification_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db NotificationRepository

type NotificationRepository interface {
	AddNotificationToken(args *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, error)
	GetAllNextMedicationsToSendNotifications() ([]models.Medication, error)
	GetSingleUserDeviceTokens(userId int) ([]string, error)
}

type notificationRepo struct {
	DB *gorm.DB
}

func NewNotificationRepo(db *GormDB) NotificationRepository {
	return &notificationRepo{db.DB}
}

func (db *notificationRepo) AddNotificationToken(args *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, error) {
	var fcmToken models.FCMNotificationToken

	fcmToken.Token = args.Token
	fcmToken.UserID = args.UserID
	err := db.DB.Create(&fcmToken).Error
	if err != nil {
		return nil, fmt.Errorf("could not create notification: %v", err)
	}

	return &fcmToken, nil
}

func (db *notificationRepo) GetAllNextMedicationsToSendNotifications() ([]models.Medication, error) {
	var medications []models.Medication

	err := db.DB.Where("date_trunc('hour', next_dosage_time) = date_trunc('hour', now())").Where("is_medication_done = false").Find(&medications).Error
	if err != nil {
		return nil, fmt.Errorf("could not get next medication: %v", err)
	}
	return medications, nil
}

func (db *notificationRepo) GetSingleUserDeviceTokens(userId int) ([]string, error) {
	var tokens []string

	err := db.DB.Table("fcm_notification_tokens").Where("user_id = ?", userId).
		Pluck("token", &tokens).Error
	if err != nil {
		return []string{}, fmt.Errorf("retrieving notification tokens: %v", err)
	}

	return tokens, nil
}
