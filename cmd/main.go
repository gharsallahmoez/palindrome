package main

import (
	"fmt"
	"github.com/gharsallahmoez/palindrome/config"
	"github.com/gharsallahmoez/palindrome/infra/database"
	"github.com/gharsallahmoez/palindrome/server/http"
	logger "github.com/sirupsen/logrus"
	"os"
)

// stop channel is used to stop the server
var stop = make(chan os.Signal)

func main() {
	conf := config.New()
	config.InitLogger()
	logger.Debug("Configuration loaded successfully")
	logger.Debug(fmt.Sprintf("Starting messages service on port %v", conf.Server.Port))

	db, err := database.Create(conf.Database)
	if err != nil {
		logger.Fatalf("failed to create the database : %v", err)
	}

	// create the service
	messageService := http.NewMessageService(db)

	srv := http.NewRunner(&conf.Server, messageService)

	// register services
	srv.RegisterServices()

	if err := srv.Start(); err != nil {
		logger.Panicf("failed to start the server : %v", err)
	}
	// schedule the stop action to wait for an os signal
	go func() {
		srv.Stop(stop)
		os.Exit(0)
	}()
}
