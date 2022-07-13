package main

import (
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
	s := &server.Server{
		Config:         conf,
		AuthRepository: authRepo,
	}
	s.Start()
}
