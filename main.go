package main

import (
	"github.com/decagonhq/meddle-api/services"
	"log"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/server"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(gormDB)
	service := services.NewAuthService(authRepo)
	s := &server.Server{
		Config:         conf,
		AuthRepository: authRepo,
		AuthService:    service,
	}
	s.Start()
}
