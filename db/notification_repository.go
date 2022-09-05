package db

import (
	"fmt"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../mocks/notification_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db NotificationRepository

type NotificationRepository interface {
	AddNotificationToken(args *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, error)
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

	err := db.DB.Create(fcmToken).Error
	if err != nil {
		return nil, fmt.Errorf("could not create medication: %v", err)
	}

	return &fcmToken, nil
}
