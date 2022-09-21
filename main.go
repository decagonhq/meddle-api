package main

import (
	"log"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/server"
	"github.com/decagonhq/meddle-api/services"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(gormDB)
	mail := services.NewMailService(conf)
	notificationRepo := db.NewNotificationRepo(gormDB)
	pushNotification, errr := services.NewFirebaseCloudMessaging(notificationRepo, conf)
	if err != nil {
		log.Fatalf("error retrieving client for push notification\n%v", errr)
	}
	authService := services.NewAuthService(authRepo, conf, mail, pushNotification)

	medicationHistoryRepo := db.NewMedicationHistoryRepo(gormDB)
	medicationRepo := db.NewMedicationRepo(gormDB)
	medicationService := services.NewMedicationService(medicationRepo, medicationHistoryRepo, conf)
	medicationHistoryService := services.NewMedicationHistoryService(medicationHistoryRepo, conf)

	s := &server.Server{
		Config:                   conf,
		AuthRepository:           authRepo,
		AuthService:              authService,
		MedicationService:        medicationService,
		MedicationHistoryService: medicationHistoryService,
		PushNotification:         pushNotification,
	}
	go services.UpdateMedicationCronJob(medicationService)
	go pushNotification.NotificationsCronJob()
	go pushNotification.NotificationsCronJobFor15MinutesEarly()
	s.Start()
}
