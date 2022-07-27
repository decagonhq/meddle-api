package server

import (
	"fmt"
	"github.com/decagonhq/meddle-api/config"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"testing"
)

var testServer struct {
	router  *gin.Engine
	handler *Server
}

func TestMain(m *testing.M) {
	fmt.Println("Starting server tests")
	c, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	testServer.handler = &Server{
		Config: c,
	}
	testServer.router = testServer.handler.setupRouter()
	exitCode := m.Run()
	os.Exit(exitCode)
}
