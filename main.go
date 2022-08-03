package main

import (
	"log"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/server"
	"github.com/decagonhq/meddle-api/services"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(gormDB)
	authService := services.NewAuthService(authRepo, conf)

	medicationRepo := db.NewMedicationRepo(gormDB)
	medicationService := services.NewMedicationService(medicationRepo, conf)
	s := &server.Server{
		Config:            conf,
		AuthRepository:    authRepo,
		AuthService:       authService,
		MedicationService: medicationService,
	}
	go db.UpdateNextMedicationCronJob(medicationRepo)
	//go db.RunCronJob(medicationRepo)
	s.Start()
}
